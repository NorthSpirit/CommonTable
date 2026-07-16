package shared

// altLanguageLabel is the display text for the alternate-language toggle
// (e.g. "LT", "ET"). Set once at startup via SetAltLanguageLabel — never
// mutated afterward, so it's safe to read from concurrent request handlers.
var altLanguageLabel = "ALT"

func SetAltLanguageLabel(label string) {
	altLanguageLabel = label
}

func AltLanguageLabel() string {
	return altLanguageLabel
}
