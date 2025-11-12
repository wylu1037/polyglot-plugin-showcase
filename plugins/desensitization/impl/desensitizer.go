// Package impl provides the implementation of the desensitizer plugin.
package impl

import (
	"errors"
	"regexp"
	"strings"
)

// DesensitzerImpl is the concrete implementation of the Desensitizer interface.
type DesensitzerImpl struct{}

// DesensitizeName desensitizes a person's name by replacing all characters except the first one with asterisks.
// Example: "张三" -> "张*", "John Doe" -> "J*** ***"
func (d *DesensitzerImpl) DesensitizeName(name string) (string, error) {
	if name == "" {
		return "", errors.New("name cannot be empty")
	}

	runes := []rune(name)
	length := len(runes)

	if length <= 1 {
		return name, nil
	}

	// Keep the first character, replace the rest with asterisks
	for i := 1; i < length; i++ {
		if runes[i] != ' ' { // Preserve spaces
			runes[i] = '*'
		}
	}

	return string(runes), nil
}

// DesensitizeTelNo desensitizes a telephone number by masking the middle 4 digits.
// Example: "13812345678" -> "138****5678"
func (d *DesensitzerImpl) DesensitizeTelNo(telNo string) (string, error) {
	if telNo == "" {
		return "", errors.New("telephone number cannot be empty")
	}

	runes := []rune(telNo)
	length := len(runes)

	// Chinese mobile phone numbers are typically 11 digits
	if length != 11 {
		return "", errors.New("invalid telephone number length, expected 11 digits")
	}

	// Mask digits from index 3 to 6 (inclusive)
	for i := 3; i <= 6; i++ {
		runes[i] = '*'
	}

	return string(runes), nil
}

// DesensitizeIDNumber desensitizes an ID card number by masking most of the middle digits.
// Example: "110101199001011234" -> "11**************34"
func (d *DesensitzerImpl) DesensitizeIDNumber(idNumber string) (string, error) {
	if idNumber == "" {
		return "", errors.New("ID number cannot be empty")
	}

	runes := []rune(idNumber)
	length := len(runes)

	// Chinese ID card numbers are typically 18 digits
	if length != 18 {
		return "", errors.New("invalid ID number length, expected 18 characters")
	}

	// Keep first 2 and last 2 characters, mask the rest
	for i := 2; i < length-2; i++ {
		runes[i] = '*'
	}

	return string(runes), nil
}

// DesensitizeEmail desensitizes an email address by masking most of the username part.
// Example: "user@example.com" -> "u***@example.com"
func (d *DesensitzerImpl) DesensitizeEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("invalid email format")
	}

	username := parts[0]
	domain := parts[1]

	runes := []rune(username)
	length := len(runes)

	if length <= 1 {
		return email, nil
	}

	// Keep the first character, mask the rest of the username
	for i := 1; i < length; i++ {
		runes[i] = '*'
	}

	return string(runes) + "@" + domain, nil
}

// DesensitizeBankCard desensitizes a bank card number by masking most of the middle digits.
// Example: "6222021234567890123" -> "622202***********123"
func (d *DesensitzerImpl) DesensitizeBankCard(cardNumber string) (string, error) {
	if cardNumber == "" {
		return "", errors.New("bank card number cannot be empty")
	}

	runes := []rune(cardNumber)
	length := len(runes)

	// Bank card numbers are typically 16-19 digits
	if length < 16 || length > 19 {
		return "", errors.New("invalid bank card number length, expected 16-19 digits")
	}

	// Keep first 6 and last 3 digits, mask the rest
	for i := 6; i < length-3; i++ {
		runes[i] = '*'
	}

	return string(runes), nil
}

// DesensitizeAddress desensitizes an address by keeping only the province and city.
// Example: "北京市朝阳区某某街道123号" -> "北京市朝阳区******"
func (d *DesensitzerImpl) DesensitizeAddress(address string) (string, error) {
	if address == "" {
		return "", errors.New("address cannot be empty")
	}

	runes := []rune(address)
	length := len(runes)

	// If address is too short, return as is
	if length <= 6 {
		return address, nil
	}

	// Keep approximately the first 1/3 of the address, mask the rest
	keepLength := length / 3
	if keepLength < 6 {
		keepLength = 6
	}

	// Mask the rest
	for i := keepLength; i < length; i++ {
		runes[i] = '*'
	}

	return string(runes), nil
}
