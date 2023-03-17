package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func GenerateUid() (string, error) {

	uid, err := uuid.NewUUID()

	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to generate UUID"))
	}

	res := uid.String()

	return res, err

}

func Hash256String(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
