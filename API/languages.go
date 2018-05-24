package main

import (
	"github.com/pkg/errors"
	http "net/http"
	"strings"
)

var (
	//AvailableLanguages : List of all available languages
	AvailableLanguages = []string{"EN", "NL", "DE", "RO"}
)

func parseLanguage(r *http.Request) (string, error) {
	lang := r.URL.Query().Get("language")
	if lang == "" {
		return "EN", nil
	}
	if isCorrectLanguage(lang) {
		return lang, nil
	}
	return lang, errors.New("Invalid Language '" + lang + "', must be one of [" + strings.Join(AvailableLanguages, ", ") + "]")
}

func isCorrectLanguage(lang string) bool {
	for _, availLang := range AvailableLanguages {
		if lang == availLang {
			return true
		}
	}
	return false
}
