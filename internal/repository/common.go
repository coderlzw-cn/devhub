package repository

import (
	"fmt"
	"math"

	"github.com/google/uuid"
)

func newUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("new uuid: %w", err)
	}
	return id.String(), nil
}

func parseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func nullableUUID(value *string) (uuid.NullUUID, error) {
	if value == nil {
		return uuid.NullUUID{}, nil
	}
	id, err := parseUUID(*value)
	if err != nil {
		return uuid.NullUUID{}, err
	}
	return uuid.NullUUID{UUID: id, Valid: true}, nil
}

func toInt32(value int) (int32, error) {
	if value < math.MinInt32 || value > math.MaxInt32 {
		return 0, fmt.Errorf("%d exceeds PostgreSQL INTEGER range", value)
	}
	return int32(value), nil
}
