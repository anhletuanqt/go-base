package helper

import (
	"encoding/json"
	"fmt"
)

func FormatJson(data interface{}, print bool) ([]byte, error) {
	bs, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	if print {
		fmt.Println(string(bs))
	}

	return bs, nil
}
