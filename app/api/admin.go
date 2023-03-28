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

func AdminGetWithdrawals(token string) ([]models.Withdrawal, error) {
	var w []models.Withdrawal
	resp, err := backendRequest(token, "GET", "/admin/withdrawal", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}

func AdminGetPayments(token string) ([]models.Payment, error) {
	var w []models.Payment
	resp, err := backendRequest(token, "GET", "/admin/payment", nil)
	if err != nil {
		return w, err
	}
	err = json.Unmarshal(resp, &w)
	return w, err
}

func AdminEditMerchant(token, id string, conf interface{}) error {
	_, err := backendRequest(token, "POST", "/admin/merchant/" + id, conf)
	return err
}
