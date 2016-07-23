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
		CaseSensitiveMatch: false,
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
