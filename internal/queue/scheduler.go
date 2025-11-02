package queue

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kubilayrn/ChronoGo/internal/model"
	"github.com/kubilayrn/ChronoGo/internal/redis"
	"github.com/kubilayrn/ChronoGo/internal/repository"
	"github.com/kubilayrn/ChronoGo/internal/sender"
)

type Scheduler struct {
	mu            sync.RWMutex
	isRunning     bool
	stopChan      chan struct{}
	ticker        *time.Ticker
	repo          *repository.MessageRepository
	webhookSender *sender.WebhookSender
	ctx           context.Context
	cancel        context.CancelFunc
	interval      time.Duration
	messageLimit  int
}

func NewScheduler(repo *repository.MessageRepository, webhookSender *sender.WebhookSender) *Scheduler {
	_ = godotenv.Load()

	intervalMinutes := getEnvAsInt("SCHEDULER_INTERVAL_MINUTES", 2)
	messageLimit := getEnvAsInt("SCHEDULER_MESSAGE_LIMIT", 2)

	return &Scheduler{
		stopChan:      make(chan struct{}),
		repo:          repo,
		webhookSender: webhookSender,
		interval:      time.Duration(intervalMinutes) * time.Minute,
		messageLimit:  messageLimit,
	}
}

func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return nil
	}

	s.isRunning = true
	s.ticker = time.NewTicker(s.interval)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	log.Printf("Scheduler started - will send %d messages every %v", s.messageLimit, s.interval)

	go s.run()

	go s.processMessages()

	return nil
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	s.isRunning = false
	if s.ticker != nil {
		s.ticker.Stop()
	}
	if s.cancel != nil {
		s.cancel()
	}

	log.Println("Scheduler stopped")
}

func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-s.ticker.C:
			s.processMessages()
		}
	}
}

func (s *Scheduler) processMessages() {
	ctx := context.Background()

	messages, err := s.repo.GetUnsentMessages(ctx, s.messageLimit)
	if err != nil {
		log.Printf("Failed to fetch unsent messages: %v", err)
		return
	}

	if len(messages) == 0 {
		log.Println("No unsent messages found")
		return
	}

	log.Printf("Processing %d messages", len(messages))

	for _, msg := range messages {
		if err := s.sendMessage(ctx, msg); err != nil {
			log.Printf("Failed to send message ID %d: %v", msg.ID, err)
		}
	}
}

func (s *Scheduler) sendMessage(ctx context.Context, msg model.Message) error {
	messageID, err := s.webhookSender.SendMessage(msg.To, msg.Content)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.repo.UpdateMessageStatus(ctx, msg.ID, model.StatusSent, messageID, &now)
	if err != nil {
		return err
	}

	if redis.Client != nil {
		if cacheErr := redis.CacheMessage(ctx, *messageID, now); cacheErr != nil {
			log.Printf("Failed to cache message to Redis: %v", cacheErr)
		} else {
			log.Printf("Cached messageId %s to Redis", messageID.String())
		}
	}

	log.Printf("Successfully sent message ID %d to %s (messageId: %s)", msg.ID, msg.To, messageID.String())
	return nil
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid value for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return intValue
}
