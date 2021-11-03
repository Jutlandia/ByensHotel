package forms_test

import (
	"testing"

	"github.com/Jutlandia/ByensHotel/internal/forms"
)

func TestLoginFormIsValid(t *testing.T) {
	formInputs := []struct {
		Username string
		Password string
		Valid    bool
	}{
		{Username: "alice", Password: "123123", Valid: true},
		{Username: "alice", Password: "", Valid: false},
		{Username: "", Password: "123123", Valid: false},
		{Username: "", Password: "", Valid: false},
	}
	form := &forms.LoginForm{}
	for _, input := range formInputs {
		form.Username = input.Username
		form.Password = input.Password
		if form.IsValid() != input.Valid {
			t.Errorf("Test failed with:\n\tUsername: %s\n\tPassword: %s\n",
				input.Username, input.Password)
		}
	}
}

func TestRegisterFormIsValid(t *testing.T) {
	formInputs := []struct {
		Username        string
		Email           string
		Password        string
		ConfirmPassword string
		Valid           bool
	}{
		{
			Username:        "alice",
			Email:           "test@mail.com",
			Password:        "123123",
			ConfirmPassword: "123123",
			Valid:           true,
		},
		{
			Username:        "",
			Email:           "test@mail.com",
			Password:        "123123",
			ConfirmPassword: "123123",
			Valid:           false,
		},
		{
			Username:        "alice",
			Email:           "",
			Password:        "123123",
			ConfirmPassword: "123123",
			Valid:           false,
		},
		{
			Username:        "alice",
			Email:           "test@mail.com",
			Password:        "",
			ConfirmPassword: "123123",
			Valid:           false,
		},
		{
			Username:        "alice",
			Email:           "test@mail.com",
			Password:        "123123",
			ConfirmPassword: "",
			Valid:           false,
		},
		{
			Username:        "alice",
			Email:           "test@mail.com",
			Password:        "123123",
			ConfirmPassword: "321321",
			Valid:           false,
		},
		{
			Username:        "",
			Email:           "",
			Password:        "",
			ConfirmPassword: "",
			Valid:           false,
		},
	}
	form := &forms.RegisterForm{}
	for _, input := range formInputs {
		form.Username = input.Username
		form.Email = input.Email
		form.Password = input.Password
		form.ConfirmPassword = input.ConfirmPassword
		if form.IsValid() != input.Valid {
			t.Errorf("Test failed with:\n\tUsername: %s\n\tEmail: %s\n\tPassword: %s\n\tConfirmPassword: %s\n",
				input.Username, input.Email, input.Password, input.ConfirmPassword)
		}
	}
}
