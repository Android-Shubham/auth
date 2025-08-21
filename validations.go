package main

import (
	"fmt"
	"regexp"
)

func validate(name string,email string, password string) error {
	if err := validateName(name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}
	if err := validateEmail(email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("name must be at least 3 characters long")
	}
	return nil
}

func isValidEmail(email string) bool {
    // Simple regex for demonstration; adjust for stricter validation if needed
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func validateEmail(email string) error {
    if len(email) < 5 || !isValidEmail(email) {
        return fmt.Errorf("invalid email format")
    }
    return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}