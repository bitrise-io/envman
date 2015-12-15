package parseutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/pointers"
)

// ParseBool ...
func ParseBool(userInputStr string) (bool, error) {
	if userInputStr == "" {
		return false, errors.New("No string to parse")
	}
	userInputStr = strings.TrimSpace(userInputStr)

	lowercased := strings.ToLower(userInputStr)
	if lowercased == "yes" || lowercased == "y" {
		return true, nil
	}
	if lowercased == "no" || lowercased == "n" {
		return false, nil
	}
	return strconv.ParseBool(lowercased)
}

// CastToString ...
func CastToString(v interface{}) string {
	value := fmt.Sprintf("%v", v)
	return value
}

// CastToStringPtr ...
func CastToStringPtr(value interface{}) *string {
	castedValue, ok := value.(string)
	if !ok {
		castedStr := CastToString(value)
		return pointers.NewStringPtr(castedStr)
	}
	return pointers.NewStringPtr(castedValue)
}

// CastToBoolPtr ...
func CastToBoolPtr(value interface{}) (*bool, bool) {
	castedValue, ok := value.(bool)
	if !ok {
		castedStr := CastToString(value)
		if castedStr == "" {
			return nil, false
		}

		casted, err := ParseBool(castedStr)
		if err != nil {
			return nil, false
		}

		castedValue = casted
	}

	return pointers.NewBoolPtr(castedValue), true
}
