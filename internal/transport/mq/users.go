package mq

import (
	"encoding/json"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
)

func (s *Server) handleProfileCreatedMsg(data []byte) bool {
	var body auth.MsgBodyProfileCreated
	if err := json.Unmarshal(data, &body); err != nil {
		return true
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return true
	}

	if err := s.serviceUsers.AddUser(userId); err != nil {
		return false
	}

	return true
}

func (s *Server) handleFirebaseImportMsg(data []byte) bool {
	var body auth.MsgBodyProfileFirebaseImport
	if err := json.Unmarshal(data, &body); err != nil {
		return true
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return true
	}

	if err := s.serviceUsers.ImportFirebaseData(userId, body.FirebaseId); err != nil {
		return false
	}

	return true
}

func (s *Server) handleProfileDeletedMsg(data []byte) bool {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return true
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return true
	}

	if err := s.serviceUsers.DeleteUser(userId); err != nil {
		return false
	}

	return true
}
