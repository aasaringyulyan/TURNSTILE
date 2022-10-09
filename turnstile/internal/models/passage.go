package models

type PassageCheck struct {
	KeyHex    string `json:"key_hex" binding:"required"`
	Direction int64  `json:"direction" binding:"required"`
}

type PassageCheckForApi struct {
	Card        uint64 `json:"card_id"`
	EmployeeID  uint64 `json:"employee_id"`
	TurnstileID uint64 `json:"turnstile_id"`
	Direction   uint64 `json:"direction"`
}
