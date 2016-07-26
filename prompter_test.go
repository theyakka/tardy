// Tardy
//
// Created by Posse in NYC
// http://goposse.com
//
// Copyright (c) 2016 Posse Productions LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package tardy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

type TestReader struct {
	Reader *bufio.Reader
}

func NewTestReader() TestReader {
	return TestReader{
		Reader: bufio.NewReader(os.Stdin),
	}
}

func (r TestReader) ReadClearText(file *os.File) (string, error) {
	return r.Reader.ReadString('\n')
}

func (r TestReader) ReadSecureText(file *os.File) (string, error) {
	return r.Reader.ReadString('\n')
}

func (r TestReader) Reset(newString string) {
	r.Reader.Reset(strings.NewReader(newString))
}

var testReader TestReader

func testSimpleEntry(t *testing.T, prompter *Prompter, input string, expectedOutput string, required Optionality, defaultValue string, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}

	prompter.Reader.(TestReader).Reset(input)
	val, _ := prompter.Prompt(SimplePrompt("Enter a value", required, defaultValue))
	fmt.Print("\n  > input: '", loggedInput, "', output: '", val, "', expected: '", expectedOutput, "'\n\n")
	matchesExpected := val == expectedOutput
	if matchesExpected == false && checkFails == true {
		t.Error("Expected output does not match input")
	}
	return matchesExpected
}

func testSecureEntry(t *testing.T, prompter *Prompter, input string, expectedOutput string, required Optionality, defaultValue string, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}
	prompter.Reader.(TestReader).Reset(input)
	val, _ := prompter.Prompt(SimpleSecurePrompt("Enter a secure value", required, defaultValue))
	fmt.Print("\n  > input: '", loggedInput, "', output: '", val, "', expected: '", expectedOutput, "'\n\n")
	matchesExpected := val == expectedOutput
	if matchesExpected == false && checkFails == true {
		t.Error("Expected output does not match input")
	}
	return matchesExpected
}

func testSingleMatchEntry(t *testing.T, prompter *Prompter, input string, allowedValues []string, expectedOutput string, required Optionality, defaultValue string, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}
	prompter.Reader.(TestReader).Reset(input)
	prompt := SingleValuePrompt("Enter a value", "["+strings.Join(allowedValues, ", ")+"]", allowedValues, required, defaultValue)
	val, _ := prompter.Prompt(prompt)
	fmt.Print("\n  > input: '", loggedInput, "', output: '", val, "', expected: '", expectedOutput, "'\n\n")
	matchesExpected := val == expectedOutput
	if matchesExpected == false && checkFails == true {
		t.Error("Expected output does not match input")
	}
	return matchesExpected
}

func testYesNoEntry(t *testing.T, prompter *Prompter, input string, expectedOutput bool, required Optionality, defaultValue bool, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}
	prompter.Reader.(TestReader).Reset(input)
	prompt := YesNoPrompt("Enter a value", "", required, defaultValue)
	val, _ := prompter.Prompt(prompt)
	fmt.Print("\n  > input: '", loggedInput, "', output: '", val, "', expected: '", expectedOutput, "'\n\n")
	matchesExpected := val == expectedOutput
	if matchesExpected == false && checkFails == true {
		t.Error("Expected output does not match input")
	}
	return matchesExpected
}

func TestMain(t *testing.T) {

	p := NewPrompter()
	p.Reader = NewTestReader()

	emptyString := "\n"

	fmt.Print("\n\n")

	testSimpleEntry(t, &p, "test 1234\n", "test 1234", Required, "", true)
	testSimpleEntry(t, &p, "test 1234\n", "test1234", Required, "", false)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test 1234", true)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test", false)

	testSecureEntry(t, &p, "password\n", "password", Required, "", true)
	testSecureEntry(t, &p, "password\n", "test1234", Required, "", false)
	testSecureEntry(t, &p, "\n", "test 1234", NotRequired, "test 1234", true)
	testSecureEntry(t, &p, "\n", "test 1234", NotRequired, "test", false)

	testSingleMatchEntry(t, &p, "red\n", []string{"red", "Green", "YELLOW", "puRple"}, "red", Required, "", true)
	testSingleMatchEntry(t, &p, "purple\n", []string{"red", "Green", "YELLOW", "puRple"}, "puRple", Required, "", true)
	testSingleMatchEntry(t, &p, "blue\n", []string{"red", "Green", "YELLOW", "puRple"}, "red", Required, "", false)
	testSingleMatchEntry(t, &p, emptyString, []string{"red", "Green", "YELLOW", "puRple"}, "puRple", NotRequired, "puRple", true)
	testSingleMatchEntry(t, &p, emptyString, []string{"red", "Green", "YELLOW", "puRple"}, "puRple", NotRequired, "blue", false)

	testYesNoEntry(t, &p, "yes\n", true, Required, false, true)
	testYesNoEntry(t, &p, "turnip\n", false, Required, false, true)
	testYesNoEntry(t, &p, "nope\n", false, Required, true, true)
	testYesNoEntry(t, &p, "YeP\n", true, Required, true, true)
	testYesNoEntry(t, &p, "y\n", true, Required, true, true)
	testYesNoEntry(t, &p, "N\n", false, Required, true, true)

	fmt.Print("\n\n")
}
