package utils

import (
	"errors"
	"regexp"
	"strings"

	"upcycle_connect-api/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var nameRegex = regexp.MustCompile(`^[a-zA-ZÀ-ÖØ-öø-ÿ' -]{2,150}$`)

func ValidateUserCreation(user models.User) error {
	user.Email = strings.TrimSpace(user.Email)
	user.PasswordUser = strings.TrimSpace(user.PasswordUser)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	if user.Email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if user.PasswordUser == "" {
		return errors.New("password is required")
	}

	if len(user.PasswordUser) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(user.PasswordUser) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	if user.FirstName == "" {
		return errors.New("first_name is required")
	}

	if !nameRegex.MatchString(user.FirstName) {
		return errors.New("Le prénom n'accepte que des lettres, espaces, tirets et apostrophes (2 caractères minimum)")
	}

	if user.LastName == "" {
		return errors.New("last_name is required")
	}

	if !nameRegex.MatchString(user.LastName) {
		return errors.New("Le nom n'accepte que des lettres, espaces, tirets et apostrophes (2 caractères minimum)")
	}

	if user.UpcyclingScore < 0 {
		return errors.New("upcycling_score cannot be negative")
	}

	if user.Premium != 0 && user.Premium != 1 {
		return errors.New("premium must be 0 or 1")
	}

	return nil
}