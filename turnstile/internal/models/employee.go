package models

// Employee – Модель 1.0
type Employee struct {
	CardNumber uint64 `json:"card_number" db:"card_number"`
	EmployeeID uint64 `json:"employee_id" db:"employee_id"`
	Rv         uint64 `json:"rv" db:"rv"`
	IsDeleted  bool   `json:"isdeleted" db:"isdeleted"`
}
