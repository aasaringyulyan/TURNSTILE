package service

import (
	"database/sql"
	"strconv"
	"time"
	"turnstile/internal/models"
	"turnstile/pkg/logging"
)

const maxLogCount = 5

type TurnstileService struct {
	logger    logging.Logger
	scheduler Scheduler
	empRepo   EmployeeRepo
	logsRepo  LogsRepo

	turnstileId uint64
	logsCount   int
}

func New(logger logging.Logger, scheduler Scheduler, empRepo EmployeeRepo, logsRepo LogsRepo, turnstileId uint64) *TurnstileService {
	return &TurnstileService{
		logger:      logger,
		scheduler:   scheduler,
		empRepo:     empRepo,
		logsRepo:    logsRepo,
		turnstileId: turnstileId,
		logsCount:   0,
	}
}

func (t *TurnstileService) StartScheduler() error {
	return t.scheduler.Start()
}

func (t *TurnstileService) CheckHandler(check models.PassageCheck) (bool, error) {
	card, err := ConvertInt(check.KeyHex, 16, 10)
	if err != nil {
		return false, err
	}

	// Анонимный проход
	if card == 000000 {
		return true, nil
	}

	_, err = t.empRepo.GetEmployeeByCard(card)
	if err == sql.ErrNoRows {
		// ToDo проверить анонимный проход
		// Пытаемся обновить данные
		go t.scheduler.UnscheduledUpdate()
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (t *TurnstileService) LogHandler(logs models.PassageLogsLinux) error {
	logsForApi, err := t.GeneratePassageLogsForApi(logs)
	if err != nil {
		return err
	}

	// Save to db
	for _, v := range logsForApi.Logs {
		err = t.logsRepo.Save(v)
		if err != nil {
			return err
		}
	}

	t.logsCount = t.logsCount + len(logsForApi.Logs)
	if t.logsCount > maxLogCount {
		go t.scheduler.UnscheduledSendLogs()
		t.logsCount = 0
	}
	t.logger.Info("Log accepted")
	return nil
}

func (t *TurnstileService) GeneratePassageCheckForApi(check models.PassageCheck) (*models.PassageCheckForApi, error) {
	// ToDo использовать HexToInt()
	card, err := ConvertInt(check.KeyHex, 16, 10)
	if err != nil {
		return nil, err
	}

	employee, err := t.empRepo.GetEmployeeByCard(card)
	if err != nil {
		return nil, err
	}

	return &models.PassageCheckForApi{
		Card:        card,
		EmployeeID:  employee.EmployeeID,
		TurnstileID: t.turnstileId,
		Direction:   uint64(check.Direction - 1),
	}, err
}

func (t *TurnstileService) GeneratePassageLogsForApi(logs models.PassageLogsLinux) (models.PassageLogsForApi, error) {
	var passageLogs models.PassageLogsForApi

	for _, value := range logs.Logs {
		card, err := ConvertInt(value.Card, 16, 10)
		if err != nil {
			return models.PassageLogsForApi{}, err
		}

		employee, err := t.empRepo.GetEmployeeByCard(card)
		if err != nil {
			return models.PassageLogsForApi{}, err
		}

		passageLogs.Logs = append(passageLogs.Logs, models.PassageLogForApi{
			TurnstileID: t.turnstileId,
			EmployeeID:  employee.EmployeeID,
			CardID:      card,
			Direction:   uint64(value.Direction - 1),
			DateTime:    time.Unix(value.Time, 0).UTC().String(),
		})
	}

	return passageLogs, nil
}

func ConvertInt(val string, base, toBase int) (uint64, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return 0, err
	}

	rez, err := strconv.ParseUint(strconv.FormatInt(i, toBase), 10, 64)
	if err != nil {
		return 0, err
	}
	return rez, nil
}

func HexToInt(val string) (uint64, error) {
	rez, err := strconv.ParseUint(val, 16, 64)
	if err != nil {
		return 0, err
	}
	return rez, nil
}
