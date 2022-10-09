package service

import "turnstile/internal/models"

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	EmployeeRepo interface {
		Save(data models.Employee) error
		SaveSlice(data []models.Employee) error
		GetEmployeeByCard(card uint64) (models.Employee, error)
	}

	LogsRepo interface {
		Save(log models.PassageLogForApi) error
		GetAll() (models.PassageLogsForApi, error)
		DeleteAll() error
	}

	DataGripper interface {
		LoadData() error
	}

	LogSender interface {
		SendLogs() error
	}

	Scheduler interface {
		Start() error
		UnscheduledUpdate()
		UnscheduledSendLogs()
	}

	Turnstile interface {
		CheckHandler(check models.PassageCheck) (bool, error)
		LogHandler(logs models.PassageLogsLinux) error
	}
)
