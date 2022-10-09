package models

type Employee struct {
	CardNumber int64 `json:"card_number" db:"card_number"`
	EmployeeID int64 `json:"employee_id" db:"employee_id"`
	Rv         int64 `json:"rv" db:"rv"`
	IsDeleted  bool  `json:"isdeleted" db:"isdeleted"`
}
