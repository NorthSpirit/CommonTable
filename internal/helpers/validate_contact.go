package helpers

import "unicode/utf8"

func ValidateContact(email, phone string) string {
	if utf8.RuneCountInString(email) > 254 {
		return "Email must be at most 254 characters."
	}
	if utf8.RuneCountInString(phone) > 30 {
		return "Phone number must be at most 30 characters."
	}
	return ""
}
