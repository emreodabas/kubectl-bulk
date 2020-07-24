package interaction

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/rand"
	"log"
	"os"
	"testing"
)

func TestPromptValue(t *testing.T) {

	randomInput := rand.String(10)
	content := []byte(randomInput)
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

	if prompt != randomInput {
		t.Errorf("Expected Prompt result %s, but got %s ", prompt, randomInput)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}
