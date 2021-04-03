package model

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"golang.org/x/crypto/bcrypt"
)

// User is user model
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Validate validates fields of user model
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Username,
			validation.Required,
			validation.Match(regexp.MustCompile("[a-zA-Z0-9]+")),
		),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
	)
}

// HashPassword makes password field crypted
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(h)
	return nil
}

// CheckPassword checks user password correct
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

// ProtoResponse checks user password correct
func (u *User) ProtoResponse() *pb.UserResponse {
	return &pb.UserResponse{
		Email:    u.Email,
		Username: u.Username,
	}
}
