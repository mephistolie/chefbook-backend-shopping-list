package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

func (r *Repository) GetPendingMessages() ([]entity.MessageData, error) {
	var msgs []entity.MessageData

	query := fmt.Sprintf(`
			SELECT event_id, type, body
			FROM %s
			WHERE processed=false
		`, inboxTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg entity.MessageData
		err := rows.Scan(&msg.EventId, &msg.Type, &msg.Body)
		if err != nil {
			log.Warn("unable to scan message row: ", err)
			continue
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (r *Repository) AddMessage(msg entity.MessageData) error {
	query := fmt.Sprintf(`
			INSERT INTO %s (event_id, type, body)
			VALUES ($1, $2, $3)
		`, inboxTable)

	if _, err := r.db.Exec(query, msg.EventId, msg.Type, msg.Body); err != nil {
		log.Error("unable to add message to outbox: ", err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) CheckMessageProcessed(eventId uuid.UUID) error {
	query := fmt.Sprintf(`
			UPDATE %s
			SET processed=true
			WHERE event_id=$2
		`, inboxTable)

	_, err := r.db.Exec(query, eventId)
	if err != nil {
		log.Warnf("unable to check message %s as processed: %s", eventId, err)
	}
	return err
}
