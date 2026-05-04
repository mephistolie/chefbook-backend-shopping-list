package amqp

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
)

func (s *Server) handleProfileCreatedMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileCreated
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	log.Infof("adding user %s...", body.UserId)

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	return s.service.CreatePersonalShoppingList(ctx, userId, messageId)
}

func (s *Server) handleFirebaseImportMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileFirebaseImport
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	log.Infof("import firebase profile %s for user %s...", body.FirebaseId, body.UserId)
	return s.service.ImportFirebaseShoppingList(ctx, userId, body.FirebaseId, messageId)
}

func (s *Server) handleProfileDeletedMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	log.Infof("deleting user %s...", body.UserId)
	return s.service.DeleteUserShoppingLists(ctx, userId, messageId)
}
