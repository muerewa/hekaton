package utils

import (
	"os"
	"regexp"

	"github.com/muerewa/hekaton/structs"
	"gopkg.in/yaml.v2"
)

// Replace environmental variables in yaml config
func expandEnvVars(content []byte) []byte {
	// Regular expressions for $VAR and ${VAR}
	reSimple := regexp.MustCompile(`\$[a-zA-Z_][a-zA-Z0-9_]*`)
	reBraced := regexp.MustCompile(`\$\{[a-zA-Z_][a-zA-Z0-9_]*\}`)

	// Convert to string
	strContent := string(content)

	// Handling ${VAR}
	strContent = reBraced.ReplaceAllStringFunc(strContent, func(match string) string {
		varName := match[2 : len(match)-1] // extract VAR from ${VAR}
		return os.Getenv(varName)
	})

	// Handling $VAR
	strContent = reSimple.ReplaceAllStringFunc(strContent, func(match string) string {
		varName := match[1:] // extract VAR from $VAR
		return os.Getenv(varName)
	})

	return []byte(strContent)
}

func LoadConfig(path string) ([]structs.Monitor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expandedContent := expandEnvVars(data)
	var monitors []structs.Monitor
	err = yaml.Unmarshal(expandedContent, &monitors)
	return monitors, err
}
