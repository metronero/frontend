package api

import (
	"encoding/json"

	"gitlab.com/moneropay/metronero/metronero-backend/app/models"
)

func GetMerchantList(token string) ([]models.Merchant, error) {
	resp, err := backendRequest(token, "GET", "/admin/merchant", nil)
	if err != nil {
		return nil, err
	}
	var m []models.Merchant
	err = json.Unmarshal(resp, &m)
	return m, err
}

func GetMerchantById(token, merchantId string) (models.Merchant, error) {
	var m models.Merchant
	resp, err := backendRequest(token, "GET", "/admin/merchant/" + merchantId, nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func GetMerchantInfo(token string) (models.MerchantDashboardInfo, error) {
	var m models.MerchantDashboardInfo
	resp, err := backendRequest(token, "GET", "/merchant", nil)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(resp, &m)
	return m, err
}

func GetMerchantPayments(token string) ([]models.Payment, error) {
	var p []models.Payment
	resp, err := backendRequest(token, "GET", "/merchant/payment", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func GetMerchantWithdrawals(token string) ([]models.Withdrawal, error) {
	var p []models.Withdrawal
	resp, err := backendRequest(token, "GET", "/merchant/withdrawal", nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(resp, &p)
	return p, err
}

func DeleteMerchantById(token, id string) error {
	_, err := backendRequest(token, "DELETE", "/admin/merchant/" + id, nil)
	return err
}
