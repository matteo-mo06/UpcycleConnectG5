package utils

import (
	"errors"
	"regexp"
	"strings"

	"upcycle_connect-api/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var nameRegex = regexp.MustCompile(`^[a-zA-ZÀ-ÿ' -]{2,150}$`)

func ValidateUserCreation(user models.User) error {
	user.Email = strings.TrimSpace(user.Email)
	user.Password_user = strings.TrimSpace(user.Password_user)
	user.First_name = strings.TrimSpace(user.First_name)
	user.Last_name = strings.TrimSpace(user.Last_name)

	if user.Email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if user.Password_user == "" {
		return errors.New("password is required")
	}

	if len(user.Password_user) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(user.Password_user) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	if user.First_name == "" {
		return errors.New("first_name is required")
	}

	if !nameRegex.MatchString(user.First_name) {
		return errors.New("invalid first_name format")
	}

	if user.Last_name == "" {
		return errors.New("last_name is required")
	}

	if !nameRegex.MatchString(user.Last_name) {
		return errors.New("invalid last_name format")
	}

	if user.Upcycling_score < 0 {
		return errors.New("upcycling_score cannot be negative")
	}

	if user.Premium != 0 && user.Premium != 1 {
		return errors.New("premium must be 0 or 1")
	}

	return nil
}