package models

import (
	"github.com/google/uuid"
	"time"
)

type Post struct{
	Uuid uuid.UUID `json:"id" gorm:"primary_key"`
	Text string `json:"text"`
	Images []string `json:"images"`
	Videos []string `json:"videos"`
	datetimeCreated time.Time `json:"datetime_created"`
	datetimeEdited time.Time `json:"datetime_edited"`
}
