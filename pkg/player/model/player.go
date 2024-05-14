package model

type CreatePlayerReq struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
