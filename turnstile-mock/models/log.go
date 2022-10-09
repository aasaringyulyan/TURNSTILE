package models

type PassageLogsLinux struct {
	Logs []PassageLogLinux `json:"logs"`
}

type PassageLogLinux struct {
	LogId int64 `json:"logId"`
	Time  int64 `json:"time"`
	//TurnstileID int64  `json:"accessPoint"`
	Direction int64  `json:"direction"`
	Card      string `json:"keyHex"`
}
