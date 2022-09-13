package sanitize

import (
	"regexp"
	"strings"
)

const (
	countryCode      = `7`
	nonDigitsPattern = `\D+`
)

func Phone(phone string) string {
	if phone == `` {
		return ``
	}
	sanitized := regexp.MustCompile(nonDigitsPattern).ReplaceAllString(phone, ``)
	if strings.HasPrefix(sanitized, `8`) || strings.HasPrefix(sanitized, `7`) {
		sanitized = sanitized[1:]
	}
	return sanitized
}

func PhoneWithCountryCode(phone string) string {
	if phone == `` {
		return ``
	}
	sanitized := Phone(phone)
	return countryCode + sanitized
}
