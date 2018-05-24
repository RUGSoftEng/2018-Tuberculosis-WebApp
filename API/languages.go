package main

var (
	//AvailableLanguages : List of all available languages
	AvailableLanguages = [4]string{"EN", "NL", "DE", "RO"}
)

func parseLanguage(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	lang := vars["language"]
	if lang == "" {
		return "EN", nil
	}
	if isCorrectLanguage(lang) {
		return lang, nil
	}	
	return lang, errors.New("Invalid Language: '" + lang + "', must be one of: [" + strings.Join(AvailableLanguages, ", ") + "]")
}

func isCorrectLanguage(lang string) (bool) {
	for _, availLang := range AvailableLanguages {
		if lang == availLang {
			return true
		}
	}
	return false
}
