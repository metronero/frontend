package api

import (
	"encoding/json"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetInstanceInfo(token string) (*models.InstanceInfo, error) {
	resp, err := backendRequest(token, "GET", "/admin/instance")
	if err != nil {
		return nil, err
	}
	var i models.InstanceInfo
	err = json.Unmarshal(resp, &i)
	return &i, err
}
