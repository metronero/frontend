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
