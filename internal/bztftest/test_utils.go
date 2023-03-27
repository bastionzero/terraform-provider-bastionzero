package bztftest

import (
	"os"
)

// SetEnvironmentVariables sets environment variables and provides a cleanup
// function to reset the modified environment variables back to their original
// value.
//
// Source:
// https://dev.to/arxeiss/auto-reset-environment-variables-when-testing-in-go-5ec
func SetEnvironmentVariables(envs map[string]string) func() {
	originalEnvs := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}

func SurroundDoubleQuote(str string) string {
	return "\"" + str + "\""
}
