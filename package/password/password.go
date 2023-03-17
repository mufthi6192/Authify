package password

import (
	"crypto/sha256"
	"fmt"
)

func Generate(password string) (string, error) {

	startPassword := "white_chocolate_beverage"
	endPassword := "buy_at_starbuck"
	pw := startPassword + password + endPassword

	hash := sha256.New()
	hash.Write([]byte(string(pw)))

	hashedPassword := hash.Sum(nil)

	password = fmt.Sprintf("%x", hashedPassword)

	return password, nil

}
