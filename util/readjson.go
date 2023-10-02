package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJson(path, name string, data interface{}) (err error) {
	byteData, err := os.ReadFile(fmt.Sprintf("%v/%v", path, name))
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteData, data)
	return
}
