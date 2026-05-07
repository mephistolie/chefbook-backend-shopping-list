package amqp

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/dependencies/service"
	amqp "github.com/wagslane/go-rabbitmq"
	"k8s.io/utils/strings/slices"
)

const queueProfiles = "shopping_list.profiles"

var supportedMsgTypes = []string{
	auth.MsgTypeProfileCreated,
	auth.MsgTypeProfileFirebaseImport,
	auth.MsgTypeProfileDeleted,
}

type Server struct {
	conn             *amqp.Conn
	consumerProfiles *amqp.Consumer
	service          service.MQ
}

func NewServer(cfg config.Amqp, service service.MQ) (*Server, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", *cfg.User, *cfg.Password, *cfg.Host, *cfg.Port, *cfg.VHost)
	conn, err := amqp.NewConn(url)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:    conn,
		service: service,
	}, nil
}

func (s *Server) Start() error {
	var err error = nil
	s.consumerProfiles, err = amqp.NewConsumer(
		s.conn,
		queueProfiles,
		amqp.WithConsumerOptionsQueueQuorum,
		amqp.WithConsumerOptionsQueueDurable,
		amqp.WithConsumerOptionsExchangeName(auth.ExchangeProfiles),
		amqp.WithConsumerOptionsExchangeKind("fanout"),
		amqp.WithConsumerOptionsExchangeDurable,
		amqp.WithConsumerOptionsExchangeDeclare,
		amqp.WithConsumerOptionsRoutingKey(""),
	)
	if err != nil {
		return err
	}

	go func() {
		if err := s.consumerProfiles.Run(s.handleDelivery); err != nil {
			log.LogWarnError(context.Background(), log.Event{
				Event:     "mq.consumer.stopped",
				Message:   "rabbitmq consumer stopped with error",
				Component: log.ComponentAMQP,
			}, err)
		}
	}()

	return nil
}

func (s *Server) handleDelivery(delivery amqp.Delivery) amqp.Action {
	messageId, err := uuid.Parse(delivery.MessageId)
	if err != nil {
		log.LogWarn(context.Background(), log.Event{
			Event:     "mq.message.invalid_id",
			Message:   "invalid message id",
			Component: log.ComponentAMQP,
			Payload: map[string]any{
				"raw_message_id": delivery.MessageId,
				"message_type":   delivery.Type,
			},
		})
		return amqp.NackDiscard
	}

	if !slices.Contains(supportedMsgTypes, delivery.Type) {
		log.LogWarn(context.Background(), log.Event{
			Event:     "mq.message.unsupported_type",
			Message:   "unsupported message type",
			Component: log.ComponentAMQP,
			MessageID: messageId.String(),
			Payload: map[string]any{
				"message_type": delivery.Type,
			},
		})
		return amqp.NackDiscard
	}

	msg := MessageData{
		Id:   messageId,
		Type: delivery.Type,
		Body: delivery.Body,
	}
	if err = s.handleMessage(msg); err != nil {
		log.LogWarnError(context.Background(), log.Event{
			Event:     "mq.message.requeued",
			Message:   "message requeued",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
			},
		}, err)
		return amqp.NackRequeue
	}

	return amqp.Ack
}

func (s *Server) handleMessage(msg MessageData) error {
	ctx := context.Background()
	log.Log(ctx, log.Event{
		Event:     "mq.message.processing",
		Message:   "processing message",
		Component: log.ComponentAMQP,
		MessageID: msg.Id.String(),
		Payload: map[string]any{
			"message_type": msg.Type,
		},
	})
	switch msg.Type {
	case auth.MsgTypeProfileCreated:
		return s.handleProfileCreatedMsg(ctx, msg.Id, msg.Body)
	case auth.MsgTypeProfileFirebaseImport:
		return s.handleFirebaseImportMsg(ctx, msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(ctx, msg.Id, msg.Body)
	default:
		log.LogWarn(ctx, log.Event{
			Event:     "mq.message.unsupported_type",
			Message:   "got unsupported message type",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
			},
		})
		return errors.New("not implemented")
	}
}

func (s *Server) Stop() error {
	s.consumerProfiles.Close()
	return s.conn.Close()
}
