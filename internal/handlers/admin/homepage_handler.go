package admin

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/htmlcontent"
	"CommonTable/internal/services/homepageservices"
	admin "CommonTable/internal/templates/admin/homepage"
)

var homepageImageFileNames = [htmlcontent.MaxImages]string{"image_one", "image_two", "image_three", "image_four"}

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

	images, err := h.homepageService.ListImages(r.Context())
	if err != nil {
		http.Error(w, "failed to load homepage images", http.StatusInternalServerError)
		return
	}

	welcomeTextAlt := ""
	if content.WelcomeTextAlt.Valid {
		welcomeTextAlt = content.WelcomeTextAlt.String
	}

	if err := admin.HomepagePage(content.WelcomeText, welcomeTextAlt, images, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// imageChange is what one submission asks for in one image slot.
type imageChange struct {
	file   *multipart.FileHeader // a newly uploaded file, if any
	remove bool                  // empty the slot
	keep   bool                  // the slot ends up with an image
	image  htmlcontent.Image     // display settings (URL filled in when applying)
}

// Update handles POST /admin/homepage
func (h *HomepageHandler) Update(w http.ResponseWriter, r *http.Request) {
	// 10 MB max per uploaded file, held in memory during parsing.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	// Current images
	current, err := h.homepageService.ListImages(r.Context())
	if err != nil {
		http.Error(w, "failed to load current homepage images", http.StatusInternalServerError)
		return
	}

	// Everything is checked first
	changes := make([]imageChange, htmlcontent.MaxImages)
	for i := range changes {
		n := strconv.Itoa(i + 1)
		label := "Image " + n

		var file *multipart.FileHeader
		if files := r.MultipartForm.File["image_"+n]; len(files) > 0 {
			file = files[0]
		}
		if file != nil {
			if err := validateImageFile(file); err != nil {
				http.Error(w, label+": "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		hadImage := current[i].URL != ""
		remove := hadImage && file == nil && r.FormValue("remove_"+n) == "1"
		keep := file != nil || (hadImage && !remove)

		change := imageChange{file: file, remove: remove, keep: keep}

		// Settings only matter for a slot that has an image.
		if keep {
			width, msg := readPercent(r, "width_"+n, label+" width (%)", htmlcontent.DefaultWidthPercent, htmlcontent.MinWidthPercent, htmlcontent.MaxWidthPercent)
			if msg != "" {
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			position, msg := readPercent(r, "position_"+n, label+" position", htmlcontent.DefaultPositionPercent, 0, 100)
			if msg != "" {
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			shape := strings.TrimSpace(r.FormValue("shape_" + n))
			if shape == "" {
				shape = htmlcontent.DefaultShape
			}
			if !htmlcontent.ValidShape(shape) {
				http.Error(w, label+" shape is not allowed.", http.StatusBadRequest)
				return
			}
			change.image = htmlcontent.Image{WidthPercent: width, Shape: shape, PositionPercent: position}
		}

		changes[i] = change
	}

	// The editor submits HTML; strip everything except the formatting
	welcomeText := htmlcontent.Sanitize(r.FormValue("welcome_text"))
	welcomeTextAlt := htmlcontent.Sanitize(r.FormValue("welcome_text_alt"))

	if htmlcontent.IsEmpty(welcomeText) {
		http.Error(w, "welcome text is required", http.StatusBadRequest)
		return
	}
	if htmlcontent.IsEmpty(welcomeTextAlt) {
		welcomeTextAlt = ""
	}

	// Apply the image changes.
	for i, change := range changes {
		slot := i + 1

		if change.remove {
			if err := h.homepageService.RemoveImage(r.Context(), slot); err != nil {
				http.Error(w, "failed to remove image", http.StatusInternalServerError)
				return
			}
			removeUploadedImages(h.uploadsDir, "homepage", homepageImageFileNames[i])
			continue
		}
		if !change.keep {
			continue
		}

		image := change.image
		image.URL = current[i].URL
		if change.file != nil {
			url, err := saveUploadedImage(h.uploadsDir, "homepage", homepageImageFileNames[i], change.file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			image.URL = url
		}

		if err := h.homepageService.SaveImage(r.Context(), slot, image); err != nil {
			http.Error(w, "failed to save image", http.StatusInternalServerError)
			return
		}
	}

	var welcomeTextAltPtr *string
	if welcomeTextAlt != "" {
		welcomeTextAltPtr = &welcomeTextAlt
	}

	if err := h.homepageService.UpdateContent(r.Context(), welcomeText, welcomeTextAltPtr); err != nil {
		http.Error(w, "failed to update homepage content", http.StatusInternalServerError)
		return
	}

	images, err := h.homepageService.ListImages(r.Context())
	if err != nil {
		http.Error(w, "failed to load homepage images", http.StatusInternalServerError)
		return
	}

	if err := admin.HomepagePage(welcomeText, welcomeTextAlt, images, "Changes saved.").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// readPercent reads a whole number from a form field and checks it against
// [min, max]. An empty field gives the fallback. The second result is an
// error message for the user, or "" if the value is fine.
func readPercent(r *http.Request, field string, label string, fallback int, min int, max int) (int, string) {
	value := fallback
	if raw := strings.TrimSpace(r.FormValue(field)); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return 0, label + " must be a whole number."
		}
		value = n
	}
	if msg := helpers.CheckRange(label, float64(value), float64(min), float64(max)); msg != "" {
		return 0, msg
	}
	return value, ""
}
