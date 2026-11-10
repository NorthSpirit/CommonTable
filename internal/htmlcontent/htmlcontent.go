package htmlcontent

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const (
	tokenSpacer = "<p>---</p>"

	imageClass  = "rounded-lg border-2 gms-border-medium"
	spacerClass = "my-6 gms-border-medium"
)

const MaxImages = 4

const (
	DefaultWidthPercent    = 50
	DefaultPositionPercent = 50
	DefaultShape           = "original"

	MinWidthPercent = 10
	MaxWidthPercent = 100
)

type Shape struct {
	Value string
	Label string
	ratio string
}

var Shapes = []Shape{
	{Value: "original", Label: "Original (not cropped)"},
	{Value: "16:9", Label: "16:9 (wide, cropped)", ratio: "16 / 9"},
	{Value: "4:3", Label: "4:3 (cropped)", ratio: "4 / 3"},
	{Value: "1:1", Label: "1:1 (square, cropped)", ratio: "1 / 1"},
	{Value: "3:4", Label: "3:4 (tall, cropped)", ratio: "3 / 4"},
}

func ValidShape(v string) bool {
	for _, s := range Shapes {
		if s.Value == v {
			return true
		}
	}
	return false
}

func shapeRatio(v string) string {
	for _, s := range Shapes {
		if s.Value == v {
			return s.ratio
		}
	}
	return ""
}

type Image struct {
	URL             string
	WidthPercent    int // of the text column's width
	Shape           string
	PositionPercent int // 0 = left edge, 50 = centered, 100 = right edge
}

var policy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "br", "strong", "em", "u", "h1", "h2")
	return p
}()

var textOnly = bluemonday.StrictPolicy()

func Sanitize(s string) string {
	return policy.Sanitize(s)
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(html.UnescapeString(textOnly.Sanitize(s))) == ""
}

func Render(content string, images []Image) string {
	out := Sanitize(content)

	out = strings.ReplaceAll(out, tokenSpacer, `<hr class="`+spacerClass+`"/>`)

	var unplaced strings.Builder
	for i, img := range images {
		if i >= MaxImages {
			break
		}
		token := "{{image" + strconv.Itoa(i+1) + "}}"

		if strings.Contains(out, token) {
			replacement := ""
			if img.URL != "" {
				replacement = imageTag(img)
			}
			out = strings.ReplaceAll(out, token, replacement)
		} else if img.URL != "" {
			unplaced.WriteString(imageTag(img))
		}
	}

	return out + unplaced.String()
}

func imageTag(img Image) string {
	width := clamp(img.WidthPercent, MinWidthPercent, MaxWidthPercent)
	position := clamp(img.PositionPercent, 0, 100)
	left := float64(100-width) * float64(position) / 100

	var style strings.Builder
	fmt.Fprintf(&style, "display:block;width:%d%%;height:auto;margin:1rem 0 1rem %.2f%%;", width, left)
	if ratio := shapeRatio(img.Shape); ratio != "" {
		style.WriteString("aspect-ratio:" + ratio + ";object-fit:cover;")
	}

	return `<img src="` + html.EscapeString(img.URL) + `" alt="Homepage image" class="` + imageClass +
		`" style="` + style.String() + `"/>`
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
