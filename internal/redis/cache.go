package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	cacheKeyPrefix = "message:"
	cacheTTL       = 24 * time.Hour
)

type MessageCache struct {
	MessageID uuid.UUID `json:"message_id"`
	SentAt    time.Time `json:"sent_at"`
}

func CacheMessage(ctx context.Context, messageID uuid.UUID, sentAt time.Time) error {
	if Client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	cache := MessageCache{
		MessageID: messageID,
		SentAt:    sentAt,
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	key := fmt.Sprintf("%s%s", cacheKeyPrefix, messageID.String())
	err = Client.Set(ctx, key, data, cacheTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to cache message: %w", err)
	}

	return nil
}

func GetCachedMessage(ctx context.Context, messageID uuid.UUID) (*MessageCache, error) {
	if Client == nil {
		return nil, fmt.Errorf("Redis client is not initialized")
	}

	key := fmt.Sprintf("%s%s", cacheKeyPrefix, messageID.String())
	data, err := Client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cached message: %w", err)
	}

	var cache MessageCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return &cache, nil
}

