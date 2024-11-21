package hash

import "golang.org/x/crypto/bcrypt"

const (
	passCost = 10
)

// EncryptPassword возвращает зашифрованый пароль
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndPass сравнивает зашифрованный пароль с другим паролем
// в случае несовпадения возвращает ошибку, иначе - nil
func CompareHashAndPass(password, realPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(realPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}
