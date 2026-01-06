package uuid

import "github.com/google/uuid"

// uuid V4を生成する
func GenerateID() string {
	return uuid.New().String()
}

// uuid V7を生成する
func GenerateIDV7() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
