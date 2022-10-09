package models

type PassageCheck struct {
	//Type        string `json:"type"`
	KeyHex    string `json:"key_hex"`
	Direction int64  `json:"direction"`
	//AccessPoint int    `json:"access_point"`
}
