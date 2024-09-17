package dto

type ResponsePaginated[T any] struct {
	Meta PaginatedMeta `json:"meta"`
	List []T           `json:"list"` // Generic data field that can hold a slice of any type T
}

type PaginatedMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
