package sender

import (
	"database/sql"
	"turnstile/internal/models"
	"turnstile/internal/service"
	"turnstile/internal/service/infrastructure"
)

type LogSender struct {
	client     *infrastructure.LogClient
	logsRepo   service.LogsRepo
	logsBuffer *models.PassageLogsForApi
}

func New(client *infrastructure.LogClient, empRepo service.LogsRepo) *LogSender {
	return &LogSender{
		client:   client,
		logsRepo: empRepo,
	}
}

func (ls *LogSender) SendLogs() error {
	var logsForSending models.PassageLogsForApi

	if ls.logsBuffer == nil {
		// Получить все логи
		logs, err := ls.logsRepo.GetAll()
		if err != nil && err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			return err
		}
		// Положить в буфер
		ls.logsBuffer = &logs

		// Удалить все логи
		err = ls.logsRepo.DeleteAll()
		if err != nil {
			return err
		}
	}

	logsForSending = *ls.logsBuffer

	// Отправить все логи
	err := ls.client.SendLogs(logsForSending)
	if err != nil {
		return err
	}

	ls.logsBuffer = nil
	return nil
}
