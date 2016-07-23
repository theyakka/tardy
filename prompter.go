package tardy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

// Prompter - the primary prompt controller
type Prompter struct {
	Reader        *bufio.Reader
	Values        map[string]interface{}
	IndexedValues []interface{}
	TrimSpace     bool
	ColonAtEnd    bool
}

// Prompt - an individual request to get back information via a prompt
type Prompt struct {
	Message            string
	DefaultValue       interface{}
	Required           Optionality
	RetryIfNoMatch     bool
	ValueHint          string
	FailIfNoMatch      bool
	CaseSensitiveMatch bool
	ValueConverter     PromptValueConverter
	ValidationFunc     PromptValueValidator
}

// NewPrompter - creates a new prompter instance
func NewPrompter() Prompter {
	reader := bufio.NewReader(os.Stdin)
	prompter := Prompter{
		Reader:        reader,
		Values:        map[string]interface{}{},
		IndexedValues: []interface{}{},
		TrimSpace:     true,
		ColonAtEnd:    true,
	}
	return prompter
}

// Prompt - prompt for a single entry and return the provided value and validity status
func (pmt *Prompter) Prompt(prompt Prompt) (interface{}, Validity) {
	fmt.Print(pmt.formattedPromptMessage(prompt))
	readString, err := pmt.Reader.ReadString('\n')
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
	if validity == IsNotValid {
		fmt.Print("ERROR: Not a valid response.\n\n")
		return pmt.Prompt(prompt)
	}

	if prompt.ValueConverter != nil {
		finalValue = prompt.ValueConverter(&prompt, finalValue.(string))
	}
	return pmt.storeValuesAndReturn(prompt, finalValue, IsValid)
}

// Do - add a series of prompts in one go. will return an array containing a map of the
// value and validity values for each prompt (i.e.: { "value" : "..", "validity" : ".." }).
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

// internal

// storeValuesAndReturn - wrapper to avoid repeating value storage
func (pmt *Prompter) storeValuesAndReturn(prompt Prompt, value interface{}, validity Validity) (interface{}, Validity) {
	pmt.IndexedValues = append(pmt.IndexedValues, value)
	pmt.Values[prompt.Message] = value
	return value, validity
}

func (pmt *Prompter) formattedPromptMessage(prompt Prompt) string {
	colon := ""
	if pmt.ColonAtEnd {
		colon = ":"
	}
	hint := ""
	if prompt.ValueHint != "" {
		hint = " " + prompt.ValueHint
	}
	return fmt.Sprintf("%s%s%s  ", prompt.Message, hint, colon)
}

// isPositiveStringValue - returns true if the string value matches the list of positive
// string values. empty or non-matched value will return the noMatchValue.
func isPositiveStringValue(value string, noMatchValue bool) bool {
	if value == "" {
		return noMatchValue
	}
	switch strings.ToLower(value) {
	case "yes", "y", "yo", "si", "yup", "ya", "yep":
		return true
	case "no", "n", "nope", "no way", "nuh uh", "nah":
		return false
	}
	return noMatchValue
}

func mapStrings(source []string, f func(string) string) []string {
	out := make([]string, len(source))
	for idx, val := range source {
		out[idx] = f(val)
	}
	return out
}
