package storage

import (
	"crypto/tls"
	"fmt"

	"github.com/Jutlandia/ByensHotel/internal/config"
	"github.com/go-ldap/ldap"
)

type UserStorage interface {
	GetWithCredentials(string, string) (User, error)
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	Create(username string, email string, password string) error
}

type LDAPServer struct {
	config.LDAP
}

func New(cfg config.LDAP) *LDAPServer {
	return &LDAPServer{cfg}
}

func (ls LDAPServer) GetWithCredentials(username string, password string) (User, error) {
	l, err := ls.connect()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	err = l.Bind(ls.BindUsername, ls.BindPassword)
	if err != nil {
		return nil, err
	}
	sr := ldap.NewSearchRequest(
		ls.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=*)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{},
		nil,
	)
	res, err := l.Search(sr)
	if err != nil {
		return nil, err
	}
	if len(res.Entries) != 1 {
		return nil, fmt.Errorf("invalid username or password")
	}
	userdn := res.Entries[0].DN
	err = l.Bind(userdn, password)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}
	return LDAPUser{res.Entries[0]}, nil
}

func (ls LDAPServer) GetByUsername(username string) (User, error) {
	l, err := ls.connect()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	err = l.Bind(ls.BindUsername, ls.BindPassword)
	if err != nil {
		return nil, err
	}
	sr := ldap.NewSearchRequest(
		ls.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=*)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{},
		nil,
	)
	res, err := l.Search(sr)
	if err != nil {
		return nil, err
	}
	if len(res.Entries) == 1 {
		return LDAPUser{res.Entries[0]}, nil
	}
	return nil, nil
}

func (ls LDAPServer) GetByEmail(email string) (User, error) {
	l, err := ls.connect()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	err = l.Bind(ls.BindUsername, ls.BindPassword)
	if err != nil {
		return nil, err
	}
	sr := ldap.NewSearchRequest(
		ls.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=*)(mail=%s))", ldap.EscapeFilter(email)),
		[]string{},
		nil,
	)
	res, err := l.Search(sr)
	if err != nil {
		return nil, err
	}
	if len(res.Entries) == 1 {
		return LDAPUser{res.Entries[0]}, nil
	}
	return nil, nil
}

func (ls LDAPServer) Create(username string, email string, password string) error {
	l, err := ls.connect()
	if err != nil {
		return err
	}
	defer l.Close()
	err = l.Bind(ls.BindUsername, ls.BindPassword)
	if err != nil {
		return err
	}
	req := ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", username, ls.BaseDN),
		Attributes: []ldap.Attribute{
			{
				Type: "objectClass",
				Vals: []string{"inetOrgPerson"},
			},
			{
				Type: "cn",
				Vals: []string{username},
			},
			{
				Type: "sn",
				Vals: []string{username},
			},
			{
				Type: "uid",
				Vals: []string{username},
			},
			{
				Type: "ou",
				Vals: []string{"people"},
			},
			{
				Type: "mail",
				Vals: []string{email},
			},
			{
				Type: "userPassword",
				Vals: []string{password},
			},
		},
	}
	return l.Add(&req)
}

func (ls LDAPServer) connect() (*ldap.Conn, error) {
	l, err := ldap.DialURL(ls.URL)
	if err != nil {
		return nil, err
	}
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	return l, nil
}

type User interface {
	Username() string
	Email() string
}

type LDAPUser struct {
	*ldap.Entry
}

func (lu LDAPUser) Username() string {
	return lu.GetAttributeValue("uid")
}

func (lu LDAPUser) Email() string {
	return lu.GetAttributeValue("mail")
}
