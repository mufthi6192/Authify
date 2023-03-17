package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {

	var field string

	val := "ok_"
	countSep := strings.Count(val, "_")

	if countSep >= 2 {
		sep := strings.SplitN(val, "_", 2)
		field = sep[1]
	} else if countSep == 1 {
		sep := strings.Split(val, "_")
		field = sep[1]
	} else {
		panic("Failed separator can't be empty")
	}

	fmt.Println(field)

}

func TestCountBlankSpace(t *testing.T) {

	str := "  "
	totalSpace := strings.Count(str, " ")

	if totalSpace != 2 {
		t.Fatal("Fail ! Space should be two | Total Space :", totalSpace)
	}
}

func TestCountCharacter(t *testing.T) {
	str := "testingzz"
	totalSpace := strings.Count(str, " ")
	fmt.Println(len(str) - totalSpace)
}
