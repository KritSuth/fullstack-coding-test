package model

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name"          json:"name"`
	Email     string             `bson:"email"         json:"email"`
	Password  string             `bson:"password"      json:"-"`
	CreatedAt time.Time          `bson:"created_at"    json:"createdAt"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("email format is invalid")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Name == "" && r.Email == "" {
		return errors.New("at least one field (name or email) is required")
	}
	if r.Email != "" && !emailRegex.MatchString(r.Email) {
		return errors.New("email format is invalid")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(r.Password) == "" {
		return errors.New("password is required")
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}
