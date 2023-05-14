package amqp

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/transport/dependencies/service"
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
		s.handleDelivery,
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

	return nil
}

func (s *Server) handleDelivery(delivery amqp.Delivery) amqp.Action {
	messageId, err := uuid.Parse(delivery.MessageId)
	if err != nil {
		log.Warn("invalid message id: ", delivery.MessageId)
		return amqp.NackDiscard
	}

	if !slices.Contains(supportedMsgTypes, delivery.Type) {
		log.Warn("unsupported message type: ", delivery.Type)
		return amqp.NackDiscard
	}

	msg := MessageData{
		Id:   messageId,
		Type: delivery.Type,
		Body: delivery.Body,
	}
	if err = s.handleMessage(msg); err != nil {
		log.Warn("requeue message ", msg.Id)
		return amqp.NackRequeue
	}

	return amqp.Ack
}

func (s *Server) handleMessage(msg MessageData) error {
	log.Infof("processing message %s with type %s", msg.Id, msg.Type)
	switch msg.Type {
	case auth.MsgTypeProfileCreated:
		return s.handleProfileCreatedMsg(msg.Id, msg.Body)
	case auth.MsgTypeProfileFirebaseImport:
		return s.handleFirebaseImportMsg(msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(msg.Id, msg.Body)
	default:
		log.Warnf("got unsupported message type %s for message %s", msg.Type, msg.Id)
		return errors.New("not implemented")
	}
}

func (s *Server) Stop() error {
	s.consumerProfiles.Close()
	return s.conn.Close()
}
