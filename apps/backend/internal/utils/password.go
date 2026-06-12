package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of the plaintext password using the default
// cost. bcrypt automatically embeds a per-hash salt, so identical passwords
// produce different hashes.
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword reports whether the plaintext password matches the stored hash.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
