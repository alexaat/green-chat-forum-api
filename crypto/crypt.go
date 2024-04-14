package crypto

import "golang.org/x/crypto/bcrypt"

// Function recrypt password
func Encrypt(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password
	}
	return string(bytes)
}

// Function compair hash and password. Returns true if there is a match
func CompairPasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}