/*
Copyright  © 2025 PACLabs
*/

package config

import (
	"os"
	"strings"
)

type PropertyResolver struct {
	props map[string]string
}

func NewPropertyResolver() *PropertyResolver {
	return &PropertyResolver{props: make(map[string]string)}
}

// Parse -D flags from cobra's args
func (pr *PropertyResolver) ParseFlags(args []string) []string {
	remaining := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-D") {
			var kvPair string
			if arg == "-D" && i+1 < len(args) {
				// -D KEY=VALUE (space separated)
				kvPair = args[i+1]
				i++
			} else if len(arg) > 2 {
				// -DKEY=VALUE (no space)
				kvPair = arg[2:]
			}

			if kvPair != "" {
				parts := strings.SplitN(kvPair, "=", 2)
				if len(parts) == 2 {
					pr.AddProperty(parts[0], parts[1])
				}
			}
		} else {
			remaining = append(remaining, arg)
		}
	}
	return remaining
}

func (pr *PropertyResolver) AddProperty(key string, value string) {
	pr.props[key] = value
}

// Resolve with precedence: command line flags > env vars > original token
func (pr *PropertyResolver) Resolve(key string) string {
	// Check -D properties first
	if val, ok := pr.props[key]; ok {
		return val
	}
	// Fall back to environment
	if val := os.Getenv(key); val != "" {
		return val
	}
	// Keep original token
	return "${" + key + "}"
}

// Replace all tokens in a string
func (pr *PropertyResolver) ExpandProperties(tokenizedStr string) string {
	return os.Expand(tokenizedStr, pr.Resolve)
}
