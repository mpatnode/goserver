package main

import (
	"fmt"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

func stripeInit() {
	if stripe.Key == "" {
		stripe.Key = "sk_test_gjczORSbZSFaZ078UXhkaOlz007qqv4xiw"
	}
}

// StripeCreateCustomer create a strip customer and add the ID to the user object
func StripeCreateCustomer(u *User) error {
	stripeInit()

	// Could be handy to toss the "Fast" id into the description field here
	// little chicken/egg problem since we don't have that yet, but certainly solvable
	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	address := &stripe.AddressParams{
		Line1:      &u.Address.Line1,
		Line2:      &u.Address.Line2,
		City:       &u.Address.City,
		PostalCode: &u.Address.PostalCode,
		State:      &u.Address.Subdivision,
	}

	params := &stripe.CustomerParams{
		Name:    &name,
		Address: address,
		Email:   &u.Email,
	}

	c, err := customer.New(params)
	if err == nil {
		u.StripeID = c.ID
	}

	return err
}
