package handler

import (
	"github.com/kubilayrn/ChronoGo/internal/queue"
	"github.com/kubilayrn/ChronoGo/internal/repository"
)

type Handler struct {
	messageRepo *repository.MessageRepository
	scheduler   *queue.Scheduler
}

func NewHandler(messageRepo *repository.MessageRepository, scheduler *queue.Scheduler) *Handler {
	return &Handler{
		messageRepo: messageRepo,
		scheduler:   scheduler,
	}
}
