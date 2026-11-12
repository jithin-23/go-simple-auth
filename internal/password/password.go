package password

import "golang.org/x/crypto/bcrypt"

// Hash returns a bcrypt hash of the plain-text password.
func Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Matches reports whether plain is the password that produced hash.
func Matches(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
