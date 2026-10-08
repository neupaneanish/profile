package redpanda

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
	"uuid"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/valkey-io/valkey-go"
	"neupaneanish.com.np/profile/internal/utils"
)

func produce[T any](
	ctx context.Context,
	key uuid.UUID,
	topic, serviceName string,
	payload T,
	client *kgo.Client,
	logger *slog.Logger,
) {
	value, err := json.Marshal(payload)
	if err != nil {
		logger.ErrorContext(ctx, "failed to produce message", "service", serviceName, "key", key, "error", err)
		return
	}

	rKey, rKeyErr := key.MarshalText()
	if rKeyErr != nil {
		logger.ErrorContext(ctx, "failed to marshal key", "key", key.String(), "service", serviceName, "error", rKeyErr)
		return
	}

	record := &kgo.Record{
		Key:       rKey,
		Value:     value,
		Timestamp: time.Now().UTC(),
		Topic:     topic,
	}

	ctx = context.WithoutCancel(ctx)

	client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			logger.ErrorContext(
				ctx, "failed to deliver message to redpanda",
				"service", serviceName,
				"topic", topic,
				"key", key.String(),
				"error", err,
			)
			return
		}
	})
}

func RootNotificationProduce(
	ctx context.Context,
	session *utils.UserSession,
	userID uuid.UUID,
	table, method, serviceName string,
	vkClient valkey.Client,
	client *kgo.Client,
	logger *slog.Logger,
) {
	actorUsername, userUsername, err := utils.GetUsernames(
		ctx,
		session.UserID,
		userID,
		session,
		vkClient,
		logger,
	)

	if err != nil {
		logger.WarnContext(ctx, "Failed to resolve notification usernames, using placeholders", "error", err)
	}

	payload := utils.RootNotification{
		ActorID:       session.UserID,
		UserID:        userID,
		ActorUsername: actorUsername,
		UserUsername:  userUsername,
		Table:         table,
		Method:        method,
	}

	produce[utils.RootNotification](
		ctx,
		userID,
		utils.RedpandaRootNotificationTopic,
		serviceName,
		payload,
		client,
		logger,
	)
}
