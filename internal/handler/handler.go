package handler

import (
	"github.com/ifox777/metrics-alerting/internal/repository"
)

type Handler struct {
	storage repository.Storage
}

func NewHandler(storage repository.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
