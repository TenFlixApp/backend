package services

import (
	"backend/res"
	"bytes"
	"io"
	"net/http"
	"os"

	"encoding/json"
)

func GetRegisterStats() ([]map[string]interface{}, error) {
	res, err := http.Post(os.Getenv("COLLECTOR_ROUTE")+"metrics/register/aggregate", "application/json", bytes.NewBuffer([]byte(res.RegisterAggregate)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// get res body content
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
