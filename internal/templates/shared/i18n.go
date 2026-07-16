package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var translations = map[string]map[string]string{}

// LoadTranslations reads one JSON file per locale (e.g. "en.json",
// "alt.json") from dir and populates the in-memory dictionary.
// Call once at startup, before the server starts handling requests.
func LoadTranslations(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read translations dir: %w", err)
	}

	loaded := make(map[string]map[string]string)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		locale := strings.TrimSuffix(entry.Name(), ".json")

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read translation file %s: %w", entry.Name(), err)
		}

		var dict map[string]string
		if err := json.Unmarshal(data, &dict); err != nil {
			return fmt.Errorf("parse translation file %s: %w", entry.Name(), err)
		}

		loaded[locale] = dict
	}

	if _, ok := loaded["en"]; !ok {
		return fmt.Errorf("no en.json found in %s — English fallback is required", dir)
	}

	translations = loaded
	return nil
}

// T looks up a static UI string for the given locale, falling back to
// English if the locale or key is missing.
func T(locale string, key string) string {
	if dict, ok := translations[locale]; ok {
		if val, ok := dict[key]; ok {
			return val
		}
	}
	return translations["en"][key]
}
