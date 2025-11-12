package impl

import (
	"testing"
)

func TestDesensitizeName(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"中文姓名", "张三", "张*", false},
		{"英文姓名", "John Doe", "J*** ***", false},
		{"单字", "李", "李", false},
		{"空字符串", "", "", true},
		{"长姓名", "欧阳娜娜", "欧***", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DesensitizeName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDesensitizeTelNo(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"正常手机号", "13812345678", "138****5678", false},
		{"另一个手机号", "18900001111", "189****1111", false},
		{"空字符串", "", "", true},
		{"长度不对", "1234567", "", true},
		{"长度过长", "123456789012", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeTelNo(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeTelNo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DesensitizeTelNo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDesensitizeIDNumber(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"正常身份证", "110101199001011234", "11**************34", false},
		{"另一个身份证", "320102198801012345", "32**************45", false},
		{"空字符串", "", "", true},
		{"长度不对", "12345", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeIDNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeIDNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DesensitizeIDNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDesensitizeEmail(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"正常邮箱", "user@example.com", "u***@example.com", false},
		{"长用户名", "verylongusername@test.com", "v***************@test.com", false},
		{"单字符用户名", "a@test.com", "a@test.com", false},
		{"空字符串", "", "", true},
		{"无效格式", "notanemail", "", true},
		{"缺少@", "user.example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DesensitizeEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDesensitizeBankCard(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"19位银行卡", "6222021234567890123", "622202**********123", false},
		{"16位银行卡", "6222021234567890", "622202*******890", false},
		{"空字符串", "", "", true},
		{"长度太短", "123456789012345", "", true},
		{"长度太长", "12345678901234567890", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeBankCard(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeBankCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DesensitizeBankCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDesensitizeAddress(t *testing.T) {
	d := &DesensitzerImpl{}

	tests := []struct {
		name       string
		input      string
		wantErr    bool
		checkLen   bool // 是否检查长度而不是精确匹配
		shouldMask bool // 是否应该包含星号
	}{
		{"中文地址", "北京市朝阳区某某街道123号", false, true, true},
		{"短地址", "北京市", false, false, false}, // 短地址不脱敏
		{"长地址", "广东省深圳市南山区科技园南区深圳湾科技生态园10栋A座", false, true, true},
		{"空字符串", "", true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.DesensitizeAddress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DesensitizeAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkLen {
				// 检查结果长度与输入相同
				if len([]rune(result)) != len([]rune(tt.input)) {
					t.Errorf("DesensitizeAddress() result length = %v, want %v", len([]rune(result)), len([]rune(tt.input)))
				}
				// 检查是否应该包含星号
				if tt.shouldMask && result == tt.input {
					t.Errorf("DesensitizeAddress() = %v, should contain asterisks", result)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkDesensitizeName(b *testing.B) {
	d := &DesensitzerImpl{}
	for i := 0; i < b.N; i++ {
		_, _ = d.DesensitizeName("张三")
	}
}

func BenchmarkDesensitizeTelNo(b *testing.B) {
	d := &DesensitzerImpl{}
	for i := 0; i < b.N; i++ {
		_, _ = d.DesensitizeTelNo("13812345678")
	}
}

func BenchmarkDesensitizeEmail(b *testing.B) {
	d := &DesensitzerImpl{}
	for i := 0; i < b.N; i++ {
		_, _ = d.DesensitizeEmail("user@example.com")
	}
}
