package config

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"neupaneanish.com.np/profile/internal/utils"
)

type Redpanda struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewRedpanda(ctx context.Context, url, group string, logger *slog.Logger) (*Redpanda, error) {
	client, clientErr := kgo.NewClient(
		kgo.SeedBrokers(url),
		kgo.ConsumerGroup(group),
	)
	if clientErr != nil {
		return nil, clientErr
	}
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return &Redpanda{
		client: client,
		logger: logger,
	}, nil
}

func (r *Redpanda) Produce(
	ctx context.Context,
	topic, serviceName string,
	payload utils.RedpandaRootEventNotificationPayload,
) {
	value, err := json.Marshal(payload)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to produce message",
			"service", serviceName,
			"topic", topic,
			"payload", payload,
			"error", err,
		)
		return
	}
	record := &kgo.Record{
		Key:       []byte(payload.UserID.String()),
		Value:     value,
		Timestamp: time.Now().UTC(),
		Topic:     topic,
	}

	r.client.Produce(ctx, record, func(rec *kgo.Record, err error) {
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to deliver message to redpanda",
				"service", serviceName,
				"topic", rec.Topic,
				"key", string(rec.Key),
				"error", err,
			)
		}
	})
}

func (r *Redpanda) Close(ctx context.Context) error {
	if err := r.client.Flush(ctx); err != nil {
		r.logger.ErrorContext(ctx, "failed to flush redpanda records on shutdown", "error", err)
		r.client.Close()
		return err
	}

	r.client.Close()
	return nil
}
