package dto

type NewPaginatedResponse[T any] struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	List  []T `json:"list"` // Generic data field that can hold a slice of any type T
}