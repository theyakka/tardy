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

func testSingleMatchEntry(t *testing.T, prompter *Prompter, input string, allowedValues []string, expectedOutput string, required Optionality, defaultValue string, checkFails bool) bool {
	loggedInput := strings.TrimSpace(input)
	if loggedInput == "" {
		loggedInput = "<empty>"
	}
	prompter.Reader.Reset(strings.NewReader(input))
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
	prompter.Reader.Reset(strings.NewReader(input))
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

	emptyString := "\n"

	fmt.Print("\n\n")

	testSimpleEntry(t, &p, "test 1234\n", "test 1234", Required, "", true)
	testSimpleEntry(t, &p, "test 1234\n", "test1234", Required, "", false)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test 1234", true)
	testSimpleEntry(t, &p, emptyString, "test 1234", NotRequired, "test", false)

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
