// Package desensitization defines the desensitization plugin interface.
// This interface is used by both the host and the plugin.
package desensitization

// Desensitizer is the interface that we're exposing as a plugin.
// All plugins must implement this interface.
type Desensitizer interface {
	// DesensitizeName desensitizes a person's name
	// Example: "张三" -> "张**"
	DesensitizeName(name string) (string, error)

	// DesensitizeTelNo desensitizes a telephone number
	// Example: "13812345678" -> "138****5678"
	DesensitizeTelNo(telNo string) (string, error)

	// DesensitizeIDNumber desensitizes an ID card number
	// Example: "110101199001011234" -> "11**************34"
	DesensitizeIDNumber(idNumber string) (string, error)

	// DesensitizeEmail desensitizes an email address
	// Example: "user@example.com" -> "u***@example.com"
	DesensitizeEmail(email string) (string, error)

	// DesensitizeBankCard desensitizes a bank card number
	// Example: "6222021234567890123" -> "622202***********123"
	DesensitizeBankCard(cardNumber string) (string, error)

	// DesensitizeAddress desensitizes an address
	// Example: "北京市朝阳区某某街道123号" -> "北京市朝阳区******"
	DesensitizeAddress(address string) (string, error)
}
