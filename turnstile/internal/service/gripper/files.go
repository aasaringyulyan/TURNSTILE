package gripper

import (
	"encoding/json"
	"errors"
	"os"
)

func readRv(filePath string) (uint64, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		err = writeRv(filePath, 0)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	// Прочитать
	type rv struct {
		Rv uint64
	}
	var result rv
	body, err := os.ReadFile(filePath)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result.Rv, nil
}

func writeRv(filePath string, rv uint64) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	mp := make(map[string]uint64)
	mp["rv"] = rv

	jsonString, _ := json.Marshal(mp)
	_, err = file.WriteString(string(jsonString))
	if err != nil {
		return err
	}

	return nil
}
