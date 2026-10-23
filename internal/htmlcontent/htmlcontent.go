// Package htmlcontent sanitizes admin-authored homepage HTML (from the
// WYSIWYG editor) and expands the placeholder tokens it may contain.
package htmlcontent

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const (
	tokenImageOne = "{{image1}}"
	tokenImageTwo = "{{image2}}"
	tokenSpacer   = "<p>---</p>"

	imageClass  = "w-full max-w-md mx-auto h-auto rounded-lg border-2 gms-border-medium block my-4"
	spacerClass = "my-6 gms-border-medium"
)

// Only the formatting the editor toolbar can produce is allowed; everything
// else (scripts, styles, links, event handlers...) is stripped.
var policy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "br", "strong", "em", "u", "h1", "h2")
	return p
}()

var textOnly = bluemonday.StrictPolicy()

// Sanitize strips everything but the allowed formatting tags.
func Sanitize(s string) string {
	return policy.Sanitize(s)
}

// IsEmpty reports whether the HTML has no visible text (e.g. Quill's
// empty "<p><br></p>").
func IsEmpty(s string) bool {
	return strings.TrimSpace(html.UnescapeString(textOnly.Sanitize(s))) == ""
}

// Render sanitizes the stored HTML, then expands {{image1}} / {{image2}}
// into images and a lone "---" paragraph into a spacer line. An uploaded
// image the admin never placed is appended at the end, so it is never
// silently dropped.
func Render(content string, imageOneURL string, imageTwoURL string) string {
	out := Sanitize(content)

	out = strings.ReplaceAll(out, tokenSpacer, `<hr class="`+spacerClass+`"/>`)

	out, foundOne := expandImage(out, tokenImageOne, imageOneURL)
	out, foundTwo := expandImage(out, tokenImageTwo, imageTwoURL)

	if imageOneURL != "" && !foundOne {
		out += imageTag(imageOneURL)
	}
	if imageTwoURL != "" && !foundTwo {
		out += imageTag(imageTwoURL)
	}

	return out
}

func expandImage(s string, token string, url string) (string, bool) {
	if !strings.Contains(s, token) {
		return s, false
	}
	replacement := ""
	if url != "" {
		replacement = imageTag(url)
	}
	return strings.ReplaceAll(s, token, replacement), true
}

func imageTag(url string) string {
	return `<img src="` + html.EscapeString(url) + `" alt="Homepage image" class="` + imageClass + `"/>`
}
