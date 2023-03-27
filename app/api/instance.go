package api

import (
	"encoding/json"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetInstanceInfo(token string) (*models.InstanceInfo, error) {
	resp, err := backendRequest(token, "GET", "/admin/instance", nil)
	if err != nil {
		return nil, err
	}
	var i models.InstanceInfo
	err = json.Unmarshal(resp, &i)
	return &i, err
}

func UpdateInstance(token string, conf *models.Instance) error {
	_, err := backendRequest(token, "POST", "/admin/instance", conf)
	return err
}
