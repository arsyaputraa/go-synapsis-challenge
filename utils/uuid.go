package utils

import "github.com/google/uuid"

func CheckUUID(id string) (*uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &uuid, nil
}
