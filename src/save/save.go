package save

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const CurrentVersion uint8 = 1

type Save struct {
	SaveVersion    uint8  `json:"saveVersion"`
	Points         uint64 `json:"points"`
	Level          uint8  `json:"level"`
	CreatedUnix    uint64 `json:"createdUnix"`
	LastOpenedUnix uint64 `json:"lastOpenedUnix"`
	TimesClicked   uint64 `json:"timesClicked"`
	PassiveIncome  uint64 `json:"passiveIncome"`
}

// Returns a blank save file structure
func Default() Save {
	return Save{
		SaveVersion:    CurrentVersion,
		Points:         0,
		Level:          1,
		CreatedUnix:    uint64(time.Now().Unix()),
		LastOpenedUnix: uint64(time.Now().Unix()),
		TimesClicked:   0,
		PassiveIncome:  0,
	}
}

// Tries to retrieve save from given json file
func FromFile(path string) (*Save, error) {
	saveFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer saveFile.Close()

	saveBytes, err := io.ReadAll(saveFile)
	if err != nil {
		return nil, err
	}

	var save Save
	err = json.Unmarshal(saveBytes, &save)
	if err != nil {
		return nil, err
	}

	return &save, nil
}

// Creates save file with given fields
func Create(path string, save Save) error {
	saveFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer saveFile.Close()

	saveJsonBytes, err := json.MarshalIndent(save, "", " ")
	if err != nil {
		return err
	}

	_, err = saveFile.Write(saveJsonBytes)
	if err != nil {
		return err
	}

	return nil
}
