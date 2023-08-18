package main

import (
	"context"

	braintree "github.com/braintree-go/braintree-go"
)

var _bt *braintree.Braintree

func getBT() *braintree.Braintree {
	if _bt == nil {
		_bt = braintree.New(
			braintree.Sandbox,
			"vf42f4mp6hpkbpms",
			"22pbhrnnyjwnjf5c",
			"e6f6a024f7c7a9164140ab90d525962c",
		)
	}
	return _bt
}

// BTCreateCustomer create a Braintree customer and add the ID to the user object
func BTCreateCustomer(u *User) error {
	bt := getBT()
	ctx := context.Background()
	cr := braintree.CustomerRequest{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email: u.Email,
	}
	c, err := bt.Customer().Create(ctx, &cr)
	if err == nil {
		u.BraintreeID = c.Id
	}

	return err
}
