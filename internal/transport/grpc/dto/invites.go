package dto

import (
	"github.com/google/uuid"
)

func NewInvitesResponse(users []uuid.UUID) []string {
	invites := make([]string, len(users))
	for i, user := range users {
		invites[i] = user.String()
	}
	return invites
}
