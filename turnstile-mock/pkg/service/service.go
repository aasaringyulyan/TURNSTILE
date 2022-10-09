package service

import (
	"turnstile-mock/models"
	"turnstile-mock/pkg/logging"
	"turnstile-mock/pkg/repository"
)

type Passage interface {
	GeneratePassage() (models.PassageCheck, error)
	GenerateLogs(check models.PassageCheck) (models.PassageLogsLinux, error)
}

type Service struct {
	Passage
}

func NewService(logger logging.Logger, repos *repository.Repository) *Service {
	return &Service{
		Passage: NewPassageService(logger, repos.Passage),
	}
}
