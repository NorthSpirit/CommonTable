package homepageservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/htmlcontent"
	"CommonTable/internal/repository"
)

type HomepageService struct {
	queries *repository.Queries
}

func NewHomepageService(queries *repository.Queries) *HomepageService {
	return &HomepageService{queries: queries}
}

func (h *HomepageService) GetContent(ctx context.Context) (repository.HomepageContent, error) {
	content, err := h.queries.GetHomepageContent(ctx)
	if err != nil {
		return repository.HomepageContent{}, fmt.Errorf("get homepage content: %w", err)
	}

	return content, nil
}

func (h *HomepageService) GetContentForDisplay(ctx context.Context, locale string) (string, error) {
	welcomeText, err := h.queries.GetHomepageContentForDisplay(ctx, locale)
	if err != nil {
		return "", fmt.Errorf("get homepage content for display: %w", err)
	}

	return welcomeText, nil
}

func (h *HomepageService) UpdateContent(
	ctx context.Context,
	welcomeText string,
	welcomeTextAlt *string,
) error {

	err := h.queries.UpdateHomepageContent(ctx, repository.UpdateHomepageContentParams{
		WelcomeText:    welcomeText,
		WelcomeTextAlt: helpers.TextFromString(welcomeTextAlt),
	})
	if err != nil {
		return fmt.Errorf("update homepage content: %w", err)
	}

	return nil
}

// index 0 is slot 1,index 1 is slot 2, and so on...
// A slot with no image has an empty URL
func (h *HomepageService) ListImages(ctx context.Context) ([]htmlcontent.Image, error) {
	rows, err := h.queries.ListHomepageImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("list homepage images: %w", err)
	}

	images := make([]htmlcontent.Image, htmlcontent.MaxImages)
	for i := range images {
		images[i] = htmlcontent.Image{
			WidthPercent:    htmlcontent.DefaultWidthPercent,
			Shape:           htmlcontent.DefaultShape,
			PositionPercent: htmlcontent.DefaultPositionPercent,
		}
	}

	for _, row := range rows {
		if row.Slot < 1 || int(row.Slot) > len(images) {
			continue
		}
		images[row.Slot-1] = htmlcontent.Image{
			URL:             row.Url,
			WidthPercent:    int(row.WidthPercent),
			Shape:           row.Shape,
			PositionPercent: int(row.PositionPercent),
		}
	}

	return images, nil
}

// Creates or replaces the image in whichever slot.
func (h *HomepageService) SaveImage(ctx context.Context, slot int, image htmlcontent.Image) error {
	err := h.queries.UpsertHomepageImage(ctx, repository.UpsertHomepageImageParams{
		Slot:            int32(slot),
		Url:             image.URL,
		WidthPercent:    int32(image.WidthPercent),
		Shape:           image.Shape,
		PositionPercent: int32(image.PositionPercent),
	})
	if err != nil {
		return fmt.Errorf("save homepage image %d: %w", slot, err)
	}

	return nil
}

// Clears out a slot.
func (h *HomepageService) RemoveImage(ctx context.Context, slot int) error {
	if err := h.queries.DeleteHomepageImage(ctx, int32(slot)); err != nil {
		return fmt.Errorf("remove homepage image %d: %w", slot, err)
	}

	return nil
}
