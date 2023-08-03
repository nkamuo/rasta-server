package slug

import (
	"regexp"
	"strings"
)

func SlugifyWithUnderscore(input string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := re.ReplaceAllString(input, "_")
	slug = strings.Trim(slug, "_")
	return slug
}
