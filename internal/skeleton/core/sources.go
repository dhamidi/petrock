package core

import (
	"strings"
)

// ArgsSource parses CLI arguments into FormSource
type ArgsSource struct {
	args []string
	data map[string][]string
}

// NewArgsSource creates an ArgsSource from command line arguments
func NewArgsSource(args []string) *ArgsSource {
	source := &ArgsSource{
		args: args,
		data: make(map[string][]string),
	}
	source.parseArgs()
	return source
}

func (a *ArgsSource) parseArgs() {
	for i := 0; i < len(a.args); i++ {
		arg := a.args[i]

		// Handle --key=value format
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				key, value := parts[0], parts[1]
				a.data[key] = append(a.data[key], value)
			}
			continue
		}

		// Handle --key value format
		if strings.HasPrefix(arg, "--") && !strings.Contains(arg, "=") {
			key := arg[2:]
			
			// Check for boolean flags
			if strings.HasPrefix(key, "no-") {
				// Handle --no-flag format (sets flag to false)
				actualKey := key[3:]
				a.data[actualKey] = append(a.data[actualKey], "false")
			} else if i+1 < len(a.args) && !strings.HasPrefix(a.args[i+1], "--") {
				// Next argument is the value
				value := a.args[i+1]
				a.data[key] = append(a.data[key], value)
				i++ // Skip the value in next iteration
			} else {
				// Boolean flag without value (defaults to true)
				a.data[key] = append(a.data[key], "true")
			}
			continue
		}
	}
}

func (a ArgsSource) Get(key string) string {
	if values, exists := a.data[key]; exists && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (a ArgsSource) GetAll(key string) []string {
	if values, exists := a.data[key]; exists {
		return values
	}
	return nil
}

func (a ArgsSource) Keys() []string {
	keys := make([]string, 0, len(a.data))
	for k := range a.data {
		keys = append(keys, k)
	}
	return keys
}

// Convenience function for parsing CLI arguments
func ParseFromArgs(args []string, target interface{}) error {
	return DefaultParser.ParseFrom(NewArgsSource(args), target)
}
