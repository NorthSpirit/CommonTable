package externalorderservices

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidApiKey  = errors.New("invalid api key")
	ErrApiKeyNotFound = errors.New("api key not found")
)

type ApiKeyService struct {
	queries *repository.Queries
}

func NewApiKeyService(queries *repository.Queries) *ApiKeyService {
	return &ApiKeyService{queries: queries}
}

func hashKey(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

func (s *ApiKeyService) Create(ctx context.Context, restaurantID int32, label string) (string, repository.RestaurantApiKey, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", repository.RestaurantApiKey{}, fmt.Errorf("generate api key: %w", err)
	}
	plain := "ct_" + hex.EncodeToString(raw)

	row, err := s.queries.CreateApiKey(ctx, repository.CreateApiKeyParams{
		RestaurantID: restaurantID,
		KeyHash:      hashKey(plain),
		KeyPrefix:    plain[:11],
		Label:        label,
	})
	if err != nil {
		return "", repository.RestaurantApiKey{}, fmt.Errorf("create api key: %w", err)
	}

	return plain, row, nil
}

func (s *ApiKeyService) Authenticate(ctx context.Context, plain string) (repository.RestaurantApiKey, error) {
	row, err := s.queries.GetActiveApiKeyByHash(ctx, hashKey(plain))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.RestaurantApiKey{}, ErrInvalidApiKey
		}
		return repository.RestaurantApiKey{}, fmt.Errorf("find api key: %w", err)
	}

	return row, nil
}

func (s *ApiKeyService) ListByRestaurant(ctx context.Context, restaurantID int32) ([]repository.RestaurantApiKey, error) {
	keys, err := s.queries.ListApiKeysByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}

	return keys, nil
}

func (s *ApiKeyService) Revoke(ctx context.Context, id int32, restaurantID int32) error {
	rows, err := s.queries.RevokeApiKey(ctx, repository.RevokeApiKeyParams{ID: id, RestaurantID: restaurantID})
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	if rows == 0 {
		return ErrApiKeyNotFound
	}

	return nil
}
