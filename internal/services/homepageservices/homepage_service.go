package homepageservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type HomepageService struct {
	queries *repository.Queries
}

func NewHomepageService(queries *repository.Queries) *HomepageService {
	return &HomepageService{queries: queries}
}

// GetContent returns raw content in both languages — for admin editing.
func (h *HomepageService) GetContent(ctx context.Context) (repository.HomepageContent, error) {
	content, err := h.queries.GetHomepageContent(ctx)
	if err != nil {
		return repository.HomepageContent{}, fmt.Errorf("get homepage content: %w", err)
	}

	return content, nil
}

// GetContentForDisplay returns locale-resolved content — for the public homepage.
func (h *HomepageService) GetContentForDisplay(ctx context.Context, locale string) (repository.GetHomepageContentForDisplayRow, error) {
	content, err := h.queries.GetHomepageContentForDisplay(ctx, locale)
	if err != nil {
		return repository.GetHomepageContentForDisplayRow{}, fmt.Errorf("get homepage content for display: %w", err)
	}

	return content, nil
}

func (h *HomepageService) UpdateContent(
	ctx context.Context,
	welcomeText string,
	welcomeTextAlt *string,
	imageOneURL *string,
	imageTwoURL *string,
) error {

	err := h.queries.UpdateHomepageContent(ctx, repository.UpdateHomepageContentParams{
		WelcomeText:    welcomeText,
		WelcomeTextAlt: helpers.TextFromString(welcomeTextAlt),
		ImageOneUrl:    helpers.TextFromString(imageOneURL),
		ImageTwoUrl:    helpers.TextFromString(imageTwoURL),
	})
	if err != nil {
		return fmt.Errorf("update homepage content: %w", err)
	}

	return nil
}
