package admin

import (
	"net/http"

	"CommonTable/internal/services/homepageservices"
	"CommonTable/internal/templates/admin"
)

type HomepageHandler struct {
	homepageService *homepageservices.HomepageService
	uploadsDir      string
}

func NewHomepageHandler(homepageService *homepageservices.HomepageService, uploadsDir string) *HomepageHandler {
	return &HomepageHandler{homepageService: homepageService, uploadsDir: uploadsDir}
}

// Edit handles GET /admin/homepage
func (h *HomepageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	content, err := h.homepageService.GetContent(r.Context())
	if err != nil {
		http.Error(w, "failed to load homepage content", http.StatusInternalServerError)
		return
	}

	welcomeTextAlt := ""
	if content.WelcomeTextAlt.Valid {
		welcomeTextAlt = content.WelcomeTextAlt.String
	}

	imageOneURL := ""
	if content.ImageOneUrl.Valid {
		imageOneURL = content.ImageOneUrl.String
	}

	imageTwoURL := ""
	if content.ImageTwoUrl.Valid {
		imageTwoURL = content.ImageTwoUrl.String
	}

	if err := admin.HomepagePage(content.WelcomeText, welcomeTextAlt, imageOneURL, imageTwoURL, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/homepage
func (h *HomepageHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 10 MB max per uploaded file, held in memory during parsing.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	// Fetch current content first, so unreplaced images keep their
	// existing URL rather than being blanked out.
	current, err := h.homepageService.GetContent(r.Context())
	if err != nil {
		http.Error(w, "failed to load current homepage content", http.StatusInternalServerError)
		return
	}

	imageOneURL := ""
	if current.ImageOneUrl.Valid {
		imageOneURL = current.ImageOneUrl.String
	}
	imageTwoURL := ""
	if current.ImageTwoUrl.Valid {
		imageTwoURL = current.ImageTwoUrl.String
	}

	if files := r.MultipartForm.File["image_one"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "homepage", "image_one", files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imageOneURL = url
	}

	if files := r.MultipartForm.File["image_two"]; len(files) > 0 {
		url, err := saveUploadedImage(h.uploadsDir, "homepage", "image_two", files[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imageTwoURL = url
	}

	welcomeText := r.FormValue("welcome_text")
	welcomeTextAlt := r.FormValue("welcome_text_alt")

	var welcomeTextAltPtr, imageOneURLPtr, imageTwoURLPtr *string
	if welcomeTextAlt != "" {
		welcomeTextAltPtr = &welcomeTextAlt
	}
	if imageOneURL != "" {
		imageOneURLPtr = &imageOneURL
	}
	if imageTwoURL != "" {
		imageTwoURLPtr = &imageTwoURL
	}

	err = h.homepageService.UpdateContent(r.Context(), welcomeText, welcomeTextAltPtr, imageOneURLPtr, imageTwoURLPtr)
	if err != nil {
		http.Error(w, "failed to update homepage content", http.StatusInternalServerError)
		return
	}

	if err := admin.HomepagePage(welcomeText, welcomeTextAlt, imageOneURL, imageTwoURL, "Changes saved.").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}
