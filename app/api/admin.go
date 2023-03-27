package api

import (
	"encoding/json"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetAdminDashboard(token string) (*models.AdminDashboardInfo, error) {
	resp, err := backendRequest(token, "GET", "/admin", nil)
	if err != nil {
		return nil, err
	}
	var a models.AdminDashboardInfo
	err = json.Unmarshal(resp, &a)
	return &a, err
}
