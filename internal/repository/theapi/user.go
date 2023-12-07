package theapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
)

var (
	fnIOCopy = io.Copy
)

func (m Repository) Register(ctx context.Context, request model.CreateRegisterRequest) (model.CreateRegisterResponse, error) {
	url := fmt.Sprintf("%s/register", m.cfg.TheApi.BaseUrl)

	var requestBuffer bytes.Buffer
	writer := multipart.NewWriter(&requestBuffer)

	fields := map[string]string{
		"username":   request.Username,
		"password":   request.Password,
		"first_name": request.FirstName,
		"last_name":  request.LastName,
		"telephone":  request.Telephone,
		// "profile_image": "https://images.unsplash.com/photo-1648075082539-ca4a311d2afa?q=80&w=1170",
		"address":  request.Address,
		"city":     request.City,
		"province": request.Province,
		"country":  request.Country,
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile("profile_image", "image.jpg")
	if err != nil {
		return model.CreateRegisterResponse{}, err
	}
	_, err = fnIOCopy(part, request.ProfileImage)
	if err != nil {
		return model.CreateRegisterResponse{}, err
	}

	// Close the multipart writer to finalize the request body
	writer.Close()

	req, err := m.fnRequestContext(ctx, http.MethodPost, url, &requestBuffer)
	if err != nil {
		return model.CreateRegisterResponse{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.CreateRegisterResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		log.Info().Str("message", "status failed debuging").Str("data", fmt.Sprintf("%+v", resp)).Msg("")
		return model.CreateRegisterResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.CreateRegisterResponse{}, err
	}

	var response model.CreateRegisterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.CreateRegisterResponse{}, err
	}

	return response, nil
}

func (m Repository) GetProfile(ctx context.Context) (model.GetProfileResponse, error) {
	url := fmt.Sprintf("%s/profile", m.cfg.TheApi.BaseUrl)

	req, err := m.fnRequestContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.GetProfileResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("TSTKRI %s", ctx.Value(model.ContextExtApiTokenKey)))

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.GetProfileResponse{}, err
	}
	defer resp.Body.Close()

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.GetProfileResponse{}, err
	}

	var response model.GetProfileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.GetProfileResponse{}, err
	}

	return response, nil
}
