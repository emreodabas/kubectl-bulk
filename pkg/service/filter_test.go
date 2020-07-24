package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var promptMock func(val string, args ...interface{}) string

type preCheckMock struct{}

func (u preCheckMock) prompt(val string, args ...interface{}) string {
	return promptMock(val, args)
}

func TestNoneFilterSelected(t *testing.T) {

	content := []byte("app=nginx")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	os.Stdin = tmpfile
	prompt := Prompt("expectedName")

	fmt.Println(prompt)
	if prompt == "" {
		t.Errorf("prompt err")
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}
