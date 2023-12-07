package model

import "mime/multipart"

// TODO: define the response for all endpoint
// bear with me, because still need confirmation for it,
// since we always got 500 status code for /register endpoint on TheApi :)

type GetArticlesResponse struct {
}

type GetArticleDetailResponse struct {
}

type GetProfileResponse struct {
}

type CreateRegisterRequest struct {
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Telephone    string         `json:"telephone"`
	ProfileImage multipart.File `json:"profile_image"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	Province     string         `json:"province"`
	Country      string         `json:"country"`
}

type CreateRegisterResponse struct {
}

type GetTokenRequest struct {
	Username string
	Password string
}

type GetTokenResponse struct {
}

type GetRefreshTokenRequest struct {
	Refresh string
}

type GetRefreshTokenResponse struct {
}
