package models

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
