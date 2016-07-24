// Tardy
//
// Created by Posse in NYC
// http://goposse.com
//
// Copyright (c) 2016 Posse Productions LLC.
// All rights reserved.
// See the LICENSE file for licensing details and requirements.

package tardy

import "strings"

// SimplePrompt - a simple prompt to get a value with no restrictions. allows for a default.
func SimplePrompt(message string, required Optionality, defaultValue string) Prompt {
	return Prompt{
		Message:            message,
		DefaultValue:       defaultValue,
		Required:           required,
		RetryIfNoMatch:     true,
		CaseSensitiveMatch: false,
	}
}

// YesNoPrompt - prompts for a yes or no answer. the final value will be a boolean. allows for a default.
func YesNoPrompt(message string, hint string, required Optionality, defaultValue bool) Prompt {
	return Prompt{
		Message:            message,
		ValueHint:          hint,
		DefaultValue:       defaultValue,
		Required:           required,
		RetryIfNoMatch:     true,
		FailIfNoMatch:      true,
		CaseSensitiveMatch: false,
		ValidationFunc: func(prompt *Prompt, value string) (string, Validity) {
			_, validity := isValidYesOrNoValue(value)
			return value, validity
		},
		ValueConverter: func(prompt *Prompt, value string) interface{} {
			return isPositiveStringValue(value, false)
		},
	}
}

// SingleValuePrompt - prompts for an answer from within a subset of valid answers. allows for a default.
func SingleValuePrompt(message string, hint string, values []string, required Optionality, defaultValue string) Prompt {
	return Prompt{
		Message:            message,
		ValueHint:          hint,
		DefaultValue:       defaultValue,
		Required:           required,
		RetryIfNoMatch:     true,
		FailIfNoMatch:      false,
		CaseSensitiveMatch: false,
		ValidationFunc: func(prompt *Prompt, value string) (string, Validity) {
			useValue := value
			originalValues := values
			checkStrings := values
			if prompt.CaseSensitiveMatch == false {
				useValue = strings.ToLower(useValue)
				checkStrings = mapStrings(checkStrings, func(s string) string {
					return strings.ToLower(s)
				})
			}
			for idx, v := range checkStrings {
				if useValue == v {
					return originalValues[idx], IsValid
				}
			}
			validity := IsValid
			if prompt.FailIfNoMatch {
				validity = IsNotValid
			}
			return defaultValue, validity
		},
	}
}

// helpers

func isValidYesOrNoValue(value string) (bool, Validity) {
	switch strings.ToLower(value) {
	case "yes", "y", "yo", "si", "yup", "ya", "yep":
		return true, IsValid
	case "no", "n", "nope", "no way", "nuh uh", "nah":
		return false, IsValid
	}
	return false, IsNotValid
}

// isPositiveStringValue - returns true if the string value matches the list of
// positive string values. empty or non-matched value will return the
// noMatchValue.
func isPositiveStringValue(value string, noMatchValue bool) bool {
	if value == "" {
		return noMatchValue
	}
	if boolVal, validity := isValidYesOrNoValue(value); validity == IsValid {
		return boolVal
	}
	return noMatchValue
}
