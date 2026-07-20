package textmarkup

import "strings"

type InlineRun struct {
	Text      string
	Bold      bool
	Underline bool
}

type BlockKind int

const (
	BlockParagraph BlockKind = iota
	BlockHeading1
	BlockHeading2
	BlockSpacer
	BlockImage
)

type Block struct {
	Kind      BlockKind
	Inlines   []InlineRun // used by Paragraph, Heading1, Heading2
	ImageSlot int         // used by BlockImage: 1 or 2
}

// Parse turns admin-authored markup text into a sequence of structured
// blocks. Recognizes: "# " / "## " headings, "---" spacer lines,
// "{{image1}}" / "{{image2}}" image placeholders, and **bold** / __underline__
// inline formatting. Everything else is treated as plain paragraph text.
func Parse(text string) []Block {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	var blocks []Block
	var paragraphLines []string

	flushParagraph := func() {
		if len(paragraphLines) == 0 {
			return
		}
		joined := strings.Join(paragraphLines, " ")
		blocks = append(blocks, Block{Kind: BlockParagraph, Inlines: parseInline(joined)})
		paragraphLines = nil
	}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)

		switch {
		case line == "":
			flushParagraph()
		case line == "---":
			flushParagraph()
			blocks = append(blocks, Block{Kind: BlockSpacer})
		case line == "{{image1}}":
			flushParagraph()
			blocks = append(blocks, Block{Kind: BlockImage, ImageSlot: 1})
		case line == "{{image2}}":
			flushParagraph()
			blocks = append(blocks, Block{Kind: BlockImage, ImageSlot: 2})
		case strings.HasPrefix(line, "## "):
			flushParagraph()
			blocks = append(blocks, Block{Kind: BlockHeading2, Inlines: parseInline(strings.TrimPrefix(line, "## "))})
		case strings.HasPrefix(line, "# "):
			flushParagraph()
			blocks = append(blocks, Block{Kind: BlockHeading1, Inlines: parseInline(strings.TrimPrefix(line, "# "))})
		default:
			paragraphLines = append(paragraphLines, line)
		}
	}
	flushParagraph()

	return blocks
}

// EnsureImagesPlaced appends BlockImage entries for any uploaded image that
// the admin never placed via {{image1}}/{{image2}}, so an uploaded image is
// never silently dropped just because the admin didn't use the token.
func EnsureImagesPlaced(blocks []Block, hasImageOne bool, hasImageTwo bool) []Block {
	placedOne := false
	placedTwo := false

	for _, b := range blocks {
		if b.Kind == BlockImage {
			if b.ImageSlot == 1 {
				placedOne = true
			}
			if b.ImageSlot == 2 {
				placedTwo = true
			}
		}
	}

	if hasImageOne && !placedOne {
		blocks = append(blocks, Block{Kind: BlockImage, ImageSlot: 1})
	}
	if hasImageTwo && !placedTwo {
		blocks = append(blocks, Block{Kind: BlockImage, ImageSlot: 2})
	}

	return blocks
}

// parseInline scans for **bold** and __underline__ toggle markers.
// Non-nested, order-independent — both can be active on the same run.
func parseInline(s string) []InlineRun {
	var runs []InlineRun
	bold := false
	underline := false
	var buf strings.Builder

	flush := func() {
		if buf.Len() == 0 {
			return
		}
		runs = append(runs, InlineRun{Text: buf.String(), Bold: bold, Underline: underline})
		buf.Reset()
	}

	i := 0
	for i < len(s) {
		if strings.HasPrefix(s[i:], "**") {
			flush()
			bold = !bold
			i += 2
			continue
		}
		if strings.HasPrefix(s[i:], "__") {
			flush()
			underline = !underline
			i += 2
			continue
		}
		buf.WriteByte(s[i])
		i++
	}
	flush()

	return runs
}
