package utils

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost bcrypt加密成本，提高安全性（默认10，提高到12）
	BcryptCost = 12
	// MinPasswordLength 最小密码长度
	MinPasswordLength = 6
	// MaxPasswordLength 最大密码长度
	MaxPasswordLength = 128
)

// ValidatePasswordStrength 验证密码强度
// 要求：至少6位，包含字母和数字，最大128位
func ValidatePasswordStrength(password string) error {
	if len(password) < MinPasswordLength {
		return errors.New("密码长度至少6位")
	}
	if len(password) > MaxPasswordLength {
		return errors.New("密码长度不能超过128位")
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter {
		return errors.New("密码必须包含至少一个字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// HashPassword 对密码进行bcrypt哈希加密
// 使用更高的成本值提高安全性
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证密码是否匹配哈希值
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
