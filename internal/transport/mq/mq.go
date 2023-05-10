package mq

import (
	"fmt"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/random"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/repository/postgres/api"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/transport/dependencies/service"
	amqp "github.com/wagslane/go-rabbitmq"
	"k8s.io/utils/strings/slices"
	"time"
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
	inbox            api.Inbox
	serviceUsers     service.Users
}

func NewServer(cfg config.Amqp, inbox api.Inbox, serviceUsers service.Users) (*Server, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", *cfg.User, *cfg.Password, *cfg.Host, *cfg.Port, *cfg.VHost)
	conn, err := amqp.NewConn(url)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:         conn,
		inbox:        inbox,
		serviceUsers: serviceUsers,
	}, nil
}

func (s *Server) Start() error {
	var err error = nil
	s.consumerProfiles, err = amqp.NewConsumer(
		s.conn,
		s.handleDelivery,
		queueProfiles,
		amqp.WithConsumerOptionsExchangeName(auth.ExchangeProfiles),
		amqp.WithConsumerOptionsExchangeKind("fanout"),
		amqp.WithConsumerOptionsExchangeDurable,
		amqp.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	go s.observeInbox()

	return nil
}

func (s *Server) handleDelivery(delivery amqp.Delivery) amqp.Action {
	eventId, err := uuid.Parse(delivery.MessageId)
	if err != nil {
		log.Warn("invalid message id: ", delivery.MessageId)
		return amqp.NackDiscard
	}

	if !slices.Contains(supportedMsgTypes, delivery.Type) {
		log.Warn("unsupported message type: ", delivery.Type)
		return amqp.NackDiscard
	}

	msg := entity.MessageData{
		EventId: eventId,
		Type:    delivery.Type,
		Body:    delivery.Body,
	}
	if err = s.inbox.AddMessage(msg); err != nil {
		return amqp.NackRequeue
	}

	go s.handleMessage(msg)

	return amqp.Ack
}

func (s *Server) observeInbox() {
	randomizeObserveTime()
	for {
		if msgs, err := s.inbox.GetPendingMessages(); err == nil {
			for _, msg := range msgs {
				go s.handleMessage(msg)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func randomizeObserveTime() {
	time.Sleep(random.DurationSeconds(10))
}

func (s *Server) handleMessage(msg entity.MessageData) {
	handled := true
	switch msg.Type {
	case auth.MsgTypeProfileCreated:
		handled = s.handleProfileCreatedMsg(msg.Body)
	case auth.MsgTypeProfileFirebaseImport:

	case auth.MsgTypeProfileDeleted:
		handled = s.handleProfileCreatedMsg(msg.Body)
	default:
	}

	if handled {
		_ = s.inbox.CheckMessageProcessed(msg.EventId)
	}
}

func (s *Server) Stop() error {
	s.consumerProfiles.Close()
	return s.conn.Close()
}
