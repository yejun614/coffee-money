package main

import "time"

type UserDTO struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required"`
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

type UserPasswordDTO struct {
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

type CallbackGithubDTO struct {
	Code string `query:"code" validate:"required"`
}

type RequestAccessTokenGithubDTO struct {
	ClientID      string `json:"client_id"`
	ClientSecrets string `json:"client_secret"`
	Code          string `json:"code"`
}

type ResponseAccessTokenGithubDTO struct {
	AccessToken string `json:"access_token" xml:"access_token" form:"access_token"`
	Scope       string `json:"scope" xml:"scope" form:"scope"`
	TokenType   string `json:"token_type" xml:"token_type" form:"token_type"`
}

type ResponseGetAuthenticatedUserGithubDTO struct {
	Name string `json:"name"`
}

type CreateLedgerDTO struct {
	StoreName   string    `json:"store_name" xml:"store_name" form:"store_name" validate:"required"`
	Balance     int       `json:"balance" xml:"balance" form:"balance" validate:"required"`
	Description string    `json:"description" xml:"description" form:"description" validate:"required"`
}

type UpdateLedgerDTO struct {
	ID          uint      `json:"id" xml:"id" form:"id" validate:"required"`
	StoreName   string    `json:"store_name" xml:"store_name" form:"store_name" validate:"required"`
	Balance     int       `json:"balance" xml:"balance" form:"balance" validate:"required"`
	Description string    `json:"description" xml:"description" form:"description" validate:"required"`
	IsDisabled  bool      `json:"is_disabled" xml:"is_disabled" form:"is_disabled" validate:"required"`
}

type SearchLedgerDTO struct {
	StoreName        string    `query:"store_name"`
	BalanceBegin     int       `query:"balance_begin"`
	BalanceEnd       int       `query:"balance_end"`
	Description      string    `query:"description"`
	IsDisabled       bool      `query:"is_disabled"`
	Username         string    `query:"username"`
	CreatedAtBegin   time.Time `query:"created_at_begin"`
	CreatedAtEnd     time.Time `query:"created_at_end"`
	UpdatedAtBegin   time.Time `query:"updated_at_begin"`
	UpdatedAtEnd     time.Time `query:"updated_at_end"`
}
