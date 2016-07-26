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

	term "golang.org/x/crypto/ssh/terminal"
)

// Validity - whether thing is valid (for better visibility)
type Validity int

const (
	// IsNotValid - The value does not match the prompt criteria
	IsNotValid Validity = 1 << iota

	// IsValid - The value matched the prompt criteria
	IsValid
)

// Optionality - whether thing is optional or not optional (for better visibility)
type Optionality int

const (
	// Required - is not optional
	Required Optionality = 1 << iota

	// NotRequired - is optional
	NotRequired
)

// PromptValueConverter - converts a provided value to a final value type
type PromptValueConverter func(*Prompt, string) interface{}

// PromptValueValidator - checks to see if the value provided matches the criteria
type PromptValueValidator func(*Prompt, string) (string, Validity)

// PromptReader - interface that provides prompt reading capabilities
type PromptReader interface {
	// ReadClearText - read a clear text prompt
	ReadClearText(file *os.File) (string, error)

	// ReadSecureText - read a secure text prompt
	ReadSecureText(file *os.File) (string, error)
}

// StandardPromptReader - default prompt reader
type StandardPromptReader struct{}

// ReadClearText - standard clear text reader
func (pr StandardPromptReader) ReadClearText(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	return reader.ReadString('\n')
}

// ReadSecureText - standard clear text reader
func (pr StandardPromptReader) ReadSecureText(file *os.File) (string, error) {
	finalString := ""
	bytes, err := term.ReadPassword(int(file.Fd()))
	if err == nil {
		finalString = string(bytes)
	}
	return finalString, err
}

// Prompter - the primary prompt controller
type Prompter struct {
	Reader PromptReader

	// Values - map of all completed values keyed on the Prompt's Message value
	Values map[string]interface{}

	// IndexedValues - all completed values based on the order they were entered
	IndexedValues []interface{}

	// TrimSpace - should we trim leading and trailing spaces?
	TrimSpace bool

	// PromptSuffix - should we add a suffix to the end of the prompt message for
	// all messages
	PromptSuffix string
}

// Prompt - an individual request to get back information via a prompt
type Prompt struct {
	// Message - the prompt message
	Message string

	// ValueHint - textual hint to the user describing what they should enter
	ValueHint string

	// SecureEntry - whether or not what the user is typing should be visible
	SecureEntry bool

	// DefaultValue - the value that will be returned if no value is entered and the
	// entry is not required
	DefaultValue interface{}

	// Required - whether or not the entry is required
	Required Optionality

	// RetryIfNoMatch - should we re-ask for a value if no match / value was entered
	RetryIfNoMatch bool

	// FailIfNoMatch - should the entry fail if it doesn't pass validation?
	FailIfNoMatch bool

	// CaseSensitiveMatch - should we do a case sensitive check of acceptable values
	CaseSensitiveMatch bool

	// ValueConverter - logic to do conversion of the string entry to your preferred
	// output type
	ValueConverter PromptValueConverter

	// ValidationFunc - logic to validate the entry
	ValidationFunc PromptValueValidator
}

// NewPrompter - creates a new prompter instance
func NewPrompter() Prompter {
	reader := StandardPromptReader{}
	prompter := Prompter{
		Reader:        reader,
		Values:        map[string]interface{}{},
		IndexedValues: []interface{}{},
		TrimSpace:     true,
		PromptSuffix:  ":  ",
	}
	return prompter
}

// Prompt - prompt for a single entry and return the provided value and validity
// status
func (pmt *Prompter) Prompt(prompt Prompt) (interface{}, Validity) {
	fmt.Print(pmt.formattedPromptMessage(prompt))

	var readString string
	var err error
	if prompt.SecureEntry {
		readString, err = pmt.Reader.ReadSecureText(os.Stdin)
	} else {
		readString, err = pmt.Reader.ReadClearText(os.Stdin)
	}

	if err != nil {
		return pmt.storeValuesAndReturn(prompt, prompt.DefaultValue, IsNotValid)
	}
	if readString == "\n" && prompt.Required == NotRequired {
		return pmt.storeValuesAndReturn(prompt, prompt.DefaultValue, IsValid)
	}
	var finalValue interface{}
	if pmt.TrimSpace {
		finalValue = strings.TrimSpace(readString)
	} else {
		finalValue = strings.TrimRight(readString, "\n")
	}

	if finalValue == "" && prompt.Required == Required {
		fmt.Print("ERROR: You must provide a value.\n\n")
		return pmt.Prompt(prompt)
	}

	validity := IsValid
	if prompt.ValidationFunc != nil {
		finalValue, validity = prompt.ValidationFunc(&prompt, finalValue.(string))
	}
	if prompt.RetryIfNoMatch && validity == IsNotValid {
		fmt.Print("ERROR: Not a valid response.\n\n")
		return pmt.Prompt(prompt)
	}

	if prompt.ValueConverter != nil {
		finalValue = prompt.ValueConverter(&prompt, finalValue.(string))
	}
	return pmt.storeValuesAndReturn(prompt, finalValue, IsValid)
}

// Do - add a series of prompts in one go. will return an array containing a map
// of the value and validity values for each prompt (i.e.: { "value" : "..",
// "validity" : ".." }).
func (pmt *Prompter) Do(prompts ...Prompt) []map[string]interface{} {
	values := []map[string]interface{}{}
	for _, prompt := range prompts {
		value, validity := pmt.Prompt(prompt) // ignore return values. final values will be stored
		values = append(values, map[string]interface{}{
			"value":    value,
			"validity": validity,
		})
	}
	return values
}

// ClearValues - clear any currently tracked values
func (pmt *Prompter) ClearValues() {
	pmt.Values = map[string]interface{}{}
	pmt.IndexedValues = []interface{}{}
}

// internal

// storeValuesAndReturn - wrapper to avoid repeating value storage
func (pmt *Prompter) storeValuesAndReturn(prompt Prompt, value interface{}, validity Validity) (interface{}, Validity) {
	pmt.IndexedValues = append(pmt.IndexedValues, value)
	pmt.Values[prompt.Message] = value
	return value, validity
}

func (pmt *Prompter) formattedPromptMessage(prompt Prompt) string {
	suffix := ""
	if pmt.PromptSuffix != "" {
		suffix = pmt.PromptSuffix
	}
	hint := ""
	if prompt.ValueHint != "" {
		hint = " " + prompt.ValueHint
	}
	return fmt.Sprintf("%s%s%s  ", prompt.Message, hint, suffix)
}

func mapStrings(source []string, f func(string) string) []string {
	out := make([]string, len(source))
	for idx, val := range source {
		out[idx] = f(val)
	}
	return out
}
