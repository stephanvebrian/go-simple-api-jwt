package theapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
)

func (m Repository) GetToken(ctx context.Context, request model.GetTokenRequest) (model.GetTokenResponse, error) {
	url := fmt.Sprintf("%s/api/token", m.cfg.TheApi.BaseUrl)

	var requestBuffer bytes.Buffer
	writer := multipart.NewWriter(&requestBuffer)

	fields := map[string]string{
		"username": request.Username,
		"password": request.Password,
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	writer.Close()

	req, err := m.fnRequestContext(ctx, http.MethodPost, url, &requestBuffer)
	if err != nil {
		return model.GetTokenResponse{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.GetTokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Info().Str("message", "status failed debuging").Str("data", fmt.Sprintf("%+v", resp)).Msg("")
		return model.GetTokenResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.GetTokenResponse{}, err
	}

	var response model.GetTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.GetTokenResponse{}, err
	}

	return response, nil
}

func (m Repository) GetRefreshToken(ctx context.Context, request model.GetRefreshTokenRequest) (model.GetRefreshTokenResponse, error) {
	url := fmt.Sprintf("%s/api/token/refresh", m.cfg.TheApi.BaseUrl)

	var requestBuffer bytes.Buffer
	writer := multipart.NewWriter(&requestBuffer)

	fields := map[string]string{
		"refresh": request.Refresh,
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	writer.Close()

	req, err := m.fnRequestContext(ctx, http.MethodPost, url, &requestBuffer)
	if err != nil {
		return model.GetRefreshTokenResponse{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.GetRefreshTokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Info().Str("message", "status failed debuging").Str("data", fmt.Sprintf("%+v", resp)).Msg("")
		return model.GetRefreshTokenResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.GetRefreshTokenResponse{}, err
	}

	var response model.GetRefreshTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.GetRefreshTokenResponse{}, err
	}

	return response, nil
}
