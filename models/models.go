package models

import "time"

type Address struct {
	Detail   string "json detail"
	Province string "json province"
	Zipcode  int    "json zipcode"
}

type Application struct {
	Name      string    "json name"
	Address   Address   "json address"
	Birthdate time.Time "json birthdate"
}

// models/models.go
