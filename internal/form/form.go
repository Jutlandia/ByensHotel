package form

import (
	"regexp"
	"strings"
)

var rxEmail = regexp.MustCompile("^\\S+@\\S+\\.\\S+$")

type Form interface {
	IsValid() bool
	AddError(field string, err string)
}

type Register struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	Errors          map[string]string
}

func (rf *Register) IsValid() bool {
	rf.Errors = make(map[string]string)
	if strings.TrimSpace(rf.Username) == "" {
		rf.AddError("Username", "Please enter a username")
	}
	if strings.TrimSpace(rf.Email) == "" {
		rf.AddError("Email", "Please enter an email")
	} else {
		match := rxEmail.Match([]byte(rf.Email))
		if match == false {
			rf.AddError("Email", "Please enter a valid email")
		}
	}
	pw := strings.TrimSpace(rf.Password)
	cpw := strings.TrimSpace(rf.ConfirmPassword)
	if pw != "" && cpw != "" && pw != cpw {
		rf.AddError("Password", "Please confirm your password")
		rf.Password = ""
		rf.ConfirmPassword = ""
	} else {
		if pw == "" {
			rf.AddError("Password", "Please enter a password")
		}
		if cpw == "" {
			rf.AddError("ConfirmPassword", "Please confirm your password")
		}
	}
	return len(rf.Errors) == 0
}

func (rf *Register) AddError(field string, err string) {
	rf.Errors[field] = err
}

type Login struct {
	Username string
	Password string
	Errors   map[string]string
}

func (lf *Login) IsValid() bool {
	lf.Errors = make(map[string]string)
	if strings.TrimSpace(lf.Username) == "" {
		lf.AddError("Username", "Please enter your username")
	}
	if strings.TrimSpace(lf.Password) == "" {
		lf.AddError("Password", "Please enter your password")
	}
	return len(lf.Errors) == 0
}

func (lf *Login) AddError(field string, err string) {
	lf.Errors[field] = err
}
