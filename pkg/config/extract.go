package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"dario.cat/mergo"
)

var (
	productPattern = regexp.MustCompile(`(?i)product`)
	bracketPattern = regexp.MustCompile(`\[([^\]]+)\]`)
)

// ConvertStringMapToAny converts a map[string]string to a map[string]any.
func ConvertStringMapToAny(m map[string]string) map[string]any {
	result := make(map[string]any)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// FlattenMapRecursive recursively flattens a nested map[string]any into a
// single-level map. Keys are concatenated with dots to represent their original
// hierarchy ("key path").
func FlattenMapRecursive(
	input map[string]any,
	prefix string,
	output map[string]any,
) {
	for key, value := range input {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + newKey
		}
		switch v := value.(type) {
		case map[string]any:
			FlattenMapRecursive(v, newKey, output)
		case map[string]string:
			newMap := ConvertStringMapToAny(v)
			FlattenMapRecursive(newMap, newKey, output)
		default:
			output[newKey] = value
		}
	}
}

// FlattenMap flattens a given input into a single-level map.  If the input is a
// map[string]any, it calls FlattenMapRecursive to flatten it, using the provided
// prefix for keys. If the input is not a map, it treats the entire input as a
// single value and assigns it to the given prefix.
func FlattenMap(input any, prefix string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	switch config := input.(type) {
	case map[string]any:
		FlattenMapRecursive(config, prefix, output)
	case map[string]string:
		newInput := ConvertStringMapToAny(config)
		FlattenMapRecursive(newInput, prefix, output)
	default:
		output[prefix] = input
	}
	return output, nil
}

func StringToBool(a any) any {
	s, ok := a.(string)
	if !ok {
		return a
	}
	lowerS := strings.ToLower(s)
	if lowerS == "true" || lowerS == "false" {
		b, err := strconv.ParseBool(lowerS)
		if err != nil {
			return a
		}
		return b
	} else {
		return a
	}
}

// HandleKeys checks if the key contains string 'Product',
// if not, return a map, with "setting" as the key, key:value as the value
// if yes, call HandleProductKeys
func ParseKeyValues(setPatterns []string) (map[string]interface{}, error) {
	if len(setPatterns) == 0 {
		return nil, fmt.Errorf("no patterns provided")
	}

	allKeyPaths := make(map[string]interface{})

	for _, pattern := range setPatterns {
		key, value, found := strings.Cut(pattern, "=")
		if !found {
			return nil, fmt.Errorf("invalid --set format: %s (expected key=value)", pattern)
		}

		configKey, err := parseConfigKey(key, value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", key, err)
		}

		// Merge results - handle accumulation for same config types
		mergo.Merge(&allKeyPaths, configKey, mergo.WithOverride)
	}

	return allKeyPaths, nil
}

func parseConfigKey(key, value string) (map[string]interface{}, error) {
	formattedValue := StringToBool(value)
	if productPattern.MatchString(key) {
		return parseProductKey(key, formattedValue)
	}

	return map[string]interface{}{
		"setting": map[string]any{key: formattedValue},
	}, nil
}

func parseProductKey(key string, value any) (map[string]interface{}, error) {
	productPart, propertyPath, found := strings.Cut(key, ".")
	if !found {
		return nil, fmt.Errorf("invalid product key format: %s (expected Product[Name].property)", key)
	}

	name, err := extractProductName(productPart)
	if err != nil {
		return nil, fmt.Errorf("failed to extract product name from %s: %w", productPart, err)
	}

	flattedValue, err := expandMap(map[string]any{propertyPath: value})
	if err != nil {
		return nil, fmt.Errorf("failed to flatten product property %s: %w", propertyPath, err)
	}

	return map[string]interface{}{
			name: flattedValue},
		nil
}

func extractProductName(s string) (string, error) {
	matches := bracketPattern.FindStringSubmatch(s)
	if len(matches) < 2 {
		return "", fmt.Errorf("no product name found in brackets")
	}
	return matches[1], nil
}

func expandMap(input map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for k, v := range input {
		parts := strings.Split(k, ".")
		currentMap := result

		for i, part := range parts {
			if i == len(parts)-1 {
				currentMap[part] = v
			} else {
				if _, exists := currentMap[part]; !exists {
					currentMap[part] = make(map[string]interface{})
				}
				nextMap, ok := currentMap[part].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("conflict at key %s", part)
				}
				currentMap = nextMap
			}
		}
	}

	return result, nil
}
