package email

import (
	"fmt"
	"testing"
)

func TestSendMail(t *testing.T) {

	data := Mail{
		FromMail:    "verification@diselesain.my.id",
		ToMail:      "mufthi.ryanda@icloud.com",
		ToName:      "Mufthi Ryanda",
		SubjectMail: "Testing for mail sending",
		BodyMail:    "Ini adalah test baru",
	}

	err := SendMail(data)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Success send mail")
}
