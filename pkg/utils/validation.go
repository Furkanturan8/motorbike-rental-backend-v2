package utils

import (
	"regexp"
)

// Telefon numarasının geçerliliğini kontrol eden fonksiyon
func ValidatePhone(phone string) bool {
	// Basit bir telefon numarası regex'i
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	match, _ := regexp.MatchString(phoneRegex, phone)
	return match
}
