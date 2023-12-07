package theapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
)

func (m Repository) GetArticles(ctx context.Context, limit, offset int64) (model.GetArticlesResponse, error) {
	url := fmt.Sprintf("%s/article?limit=%d&offset=%d", m.cfg.TheApi.BaseUrl, limit, offset)

	req, err := m.fnRequestContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.GetArticlesResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("TSTKRI %s", ctx.Value(model.ContextExtApiTokenKey)))

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.GetArticlesResponse{}, err
	}
	defer resp.Body.Close()

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.GetArticlesResponse{}, err
	}

	var response model.GetArticlesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.GetArticlesResponse{}, err
	}

	return response, nil
}

func (m Repository) GetArticleById(ctx context.Context, id int64) (model.GetArticleDetailResponse, error) {
	url := fmt.Sprintf("%s/article/%d", m.cfg.TheApi.BaseUrl, id)

	req, err := m.fnRequestContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.GetArticleDetailResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("TSTKRI %s", ctx.Value(model.ContextExtApiTokenKey)))

	resp, err := m.fnClientDo(http.DefaultClient, req)
	if err != nil {
		return model.GetArticleDetailResponse{}, err
	}
	defer resp.Body.Close()

	body, err := m.fnParseBody(resp.Body)
	if err != nil {
		return model.GetArticleDetailResponse{}, err
	}

	var response model.GetArticleDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return model.GetArticleDetailResponse{}, err
	}

	return response, nil
}
