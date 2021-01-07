package logger

import (
	"context"
	"github.com/google/uuid"
	"github.com/hayashiki/mentions/pkg/handler"
	log "github.com/sirupsen/logrus"
)

func New(ctx context.Context) *log.Entry {
	entry := log.NewEntry(log.StandardLogger())
	if reqID, ok := ctx.Value(handler.ReqIDKey).(uuid.UUID); ok {
		entry = entry.WithField("reqID", reqID)
	}
	if userID, ok := ctx.Value(handler.UserIDKey).(uuid.UUID); ok {
		entry = entry.WithField("userID", userID)
	}
	return entry
}
