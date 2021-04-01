package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Uuid uuid.UUID `json:"id" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender string `json:"gender"`
	CurrentCity string `json:"current_city"`
	HomeTown string `json:"hometown"`
	Bio string `json:"bio"`
	DateJoined time.Time `json:"date_joined"`
	Password string `json:"password"`
	Wallet string `json:"wallet"`
	Posts []Post `json:"posts"`
	Friends []User  `json:"friends"`
}




