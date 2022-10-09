package models

type PassageLogsLinux struct {
	Logs []PassageLogLinux `json:"logs" binding:"required"`
}

type PassageLogLinux struct {
	LogId     int64  `json:"logId" binding:"required"`
	Time      int64  `json:"time" binding:"required"`
	Direction int64  `json:"direction" binding:"required"`
	Card      string `json:"keyHex" binding:"required"`
}

type PassageLogsForApi struct {
	Logs []PassageLogForApi `json:"data"`
}

type PassageLogForApi struct {
	TurnstileID uint64 `json:"turnstile_id" db:"turnstile_id"`
	EmployeeID  uint64 `json:"employee_id" db:"employee_id"`
	CardID      uint64 `json:"card" db:"card"`
	Direction   uint64 `json:"direction" db:"direction"`
	DateTime    string `json:"dt" db:"dt"`
}
