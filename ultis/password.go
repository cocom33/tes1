package ultis

import "golang.org/x/crypto/bcrypt"

func HashingPassword(password string) (string, error) {
	hashedpw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	
	return string(hashedpw), err
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	// Compare the provided password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}