package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/moneropay/metronero/metronero-frontend/utils/config"
)

func backendRequest(token, method, endpoint string, body interface{}) ([]byte, error) {
	endpoint, err := url.JoinPath(config.Backend, endpoint)
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, endpoint, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
