package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`

	Email string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
	Username string `json:"username" bson:"username" validate:"required"`
}

type Url struct {
	mgm.DefaultModel `bson:",inline"`

	LongUrl string `bson:"longurl" json:"longurl" validate:"required"`
	ShortUrl string `bson:"shorturl" json:"shorturl" validate:"required"`
	User *string `bson:"user" json:"user,omitempty"`
}
