package services

import (
	"backend/res"
	"bytes"
	"io"
	"net/http"
	"os"

	"encoding/json"
)

func unmarshalArrayFromResponse(res *http.Response, err error) ([]map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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

func GetLoginStats() ([]map[string]interface{}, error) {
	res, err := http.Post(os.Getenv("COLLECTOR_ROUTE")+"metrics/login/aggregate", "application/json", bytes.NewBuffer([]byte(res.GetLoginStats())))
	return unmarshalArrayFromResponse(res, err)
}

func GetRegisterStats() ([]map[string]interface{}, error) {
	res, err := http.Post(os.Getenv("COLLECTOR_ROUTE")+"metrics/register/aggregate", "application/json", bytes.NewBuffer([]byte(res.RegisterAggregate)))
	return unmarshalArrayFromResponse(res, err)
}

func GetUserMetrics() (adminUserCount int, disabledUserCount int, err error) {
	res, err := http.Get(os.Getenv("GUARDIAN_ROUTE") + "metrics")
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()

	type UserMetrics struct {
		AdminUserCount    int `json:"adminUserCount"`
		DisabledUserCount int `json:"disabledUserCount"`
	}

	var data UserMetrics
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return 0, 0, err
	}

	return data.AdminUserCount, data.DisabledUserCount, nil
}
