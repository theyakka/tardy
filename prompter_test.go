package tardy

import (
	"fmt"
	"strings"
	"testing"
)

func testSimpleEntry(t *testing.T, prompter *Prompter, input string, expectedOutput string, required Optionality, defaultValue string, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}
	prompter.Reader.Reset(strings.NewReader(input))
	val, _ := prompter.Prompt(SimplePrompt("Enter a value", required, defaultValue))
	fmt.Print("\n  > input: '", loggedInput, "', output: '", val, "', expected: '", expectedOutput, "'\n\n")
	matchesExpected := val == expectedOutput
	if matchesExpected == false && checkFails == true {
		t.Error("Expected output does not match input")
	}
	return matchesExpected
}

func TestMain(t *testing.T) {

	p := NewPrompter()

	emptyString := "\n"

	fmt.Print("\n\n")

	testSimpleEntry(t, &p, "test 1234\n", "test 1234", Required, "", true)
	testSimpleEntry(t, &p, "test 1234\n", "test1234", Required, "", false)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test 1234", true)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test", false)

	fmt.Print("\n\n")
}
