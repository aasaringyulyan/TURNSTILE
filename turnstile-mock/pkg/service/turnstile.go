package service

import (
	"math/rand"
	"strconv"
	"time"
	"turnstile-mock/models"
	"turnstile-mock/pkg/logging"
	"turnstile-mock/pkg/repository"
)

type PassageService struct {
	logger logging.Logger
	repo   repository.Passage
}

func NewPassageService(logger logging.Logger, repo repository.Passage) *PassageService {
	return &PassageService{
		logger: logger,
		repo:   repo,
	}
}

func (s *PassageService) GeneratePassage() (models.PassageCheck, error) {
	employees, err := s.repo.GetEmployees()
	if err != nil {
		return models.PassageCheck{}, err
	}

	rand.Seed(time.Now().Unix())

	passage := models.PassageCheck{
		KeyHex:    strconv.FormatUint(employees[rand.Intn(len(employees))].CardNumber, 16),
		Direction: 1 + rand.Int63n(2),
	}

	return passage, nil
}

func (s *PassageService) GenerateLogs(passage models.PassageCheck) (models.PassageLogsLinux, error) {
	var passageLogsLinux models.PassageLogsLinux

	passageLogsLinux.Logs = append(passageLogsLinux.Logs, models.PassageLogLinux{
		LogId:     0,
		Time:      time.Now().UTC().Unix(),
		Direction: passage.Direction,
		Card:      passage.KeyHex,
	})

	return passageLogsLinux, nil
}
