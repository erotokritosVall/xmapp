package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	usersApp "github.com/erotokritosVall/xmapp/internal/users/application"
)

type loginResponse struct {
	Token string `json:"data"`
}

func Login(ctx context.Context, cfg *TestsConfig) (*string, error) {
	req := &usersApp.LoginRequest{
		Email:    cfg.TestEmail,
		Password: cfg.TestPass,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%s/v1/login", cfg.AppHost, cfg.AppPort)

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := &loginResponse{}
	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return nil, err
	}

	return &respBody.Token, nil
}
