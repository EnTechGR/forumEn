package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// loadEnv function to manually load .env file into environment variables
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore empty lines and lines starting with #
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split the line into key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			// Set the environment variable
			err := os.Setenv(key, value)
			if err != nil {
				return fmt.Errorf("Error setting environment variable: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading .env file: %v", err)
	}

	return nil
}
