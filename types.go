package main

import (
	"math/rand"
	"time"
)


type TransferRequest struct{
	ToAccount int64 `json:"toAccount"`
	Amount int `json:"amount"`
}



type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

// this is a constructor function for creating account :
func NewAccount(firstName, lastName string) *Account {
	return &Account{
		// ID:        rand.Intn(100000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)), // if any value is not initialized then default value is assigned
		CreatedAt: time.Now().Local(),
	}
}
