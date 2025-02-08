package models

import "time"

// Estructura del cliente en MongoDB
type Customer struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	Email     string    `json:"email" bson:"email"`
	Phone     string    `json:"phone_number" bson:"phone_number"`
	Address   string    `json:"address" bson:"address"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
