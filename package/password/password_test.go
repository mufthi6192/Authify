package password

import (
	"fmt"
	"testing"
)

func TestMatchPassword(t *testing.T) {

	var err error

	password := "uniqlo"
	passwordComparer := "uniqlo"

	password, err = Generate(password)

	if err != nil {
		t.Fatal("Failed to generate password")
	}

	passwordComparer, err = Generate(passwordComparer)

	if err != nil {
		t.Fatal("Failed to generate password")
	}

	if password != passwordComparer {
		t.Fatal("Failed ! Password not match")
	}

	fmt.Println("Matched Password")
}
