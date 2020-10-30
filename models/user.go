package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FirstName string         `form:"firstName" json: "firstName" binding:"required"`
	LastName  string         `form:"lastName" json:"lastName" binding:"required"`
	Email     string         `form:"email" json:"email" gorm:"uniqueIndex" binding:"required" validate:"email"`
	Adresses  pq.StringArray `form:"addresses" json:"addresses" gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) GetFullname() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// BeforeCreate hook checks if user exists
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	var user User
	err = tx.First(&user, "email = ?", u.Email).Error
	if err == nil && u.Email == user.Email {
		return errors.New("user already exists")
	}
	return nil
}
