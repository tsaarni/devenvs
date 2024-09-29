package main

import (
	"log/slog"

	"github.com/go-ldap/ldap/v3"
)

func bind(conn *ldap.Conn, username, password string) {
	controls := []ldap.Control{
		ldap.NewControlBeheraPasswordPolicy(),
	}
	request := ldap.NewSimpleBindRequest(username, password, controls)
	response, err := conn.SimpleBind(request)

	slog.Info("Response", "error", err, "response", response)

	// Print out password policy response control
	responseControl := ldap.FindControl(response.Controls, ldap.ControlTypeBeheraPasswordPolicy)
	if responseControl == nil {
		slog.Error("Failed to find control")
		return
	}

	ppolicyResponse := responseControl.(*ldap.ControlBeheraPasswordPolicy)
	slog.Info("Password policy response", "ppolicyResponse", ppolicyResponse)


}

func passwordModify(conn *ldap.Conn, username, oldPassword, newPassword string) {

	// use password modify extended operation to change password
	request := ldap.NewPasswordModifyRequest(username, oldPassword, newPassword)
	response, err := conn.PasswordModify(request)
	if err != nil {
		slog.Error("Failed to modify password", "error", err)
		return
	}

	slog.Info("Password modify response", "passwordModifyResponse", response)
}

func main() {

	// connect to LDAP server
	conn, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		slog.Error("Failed to dial", "error", err)
		return
	}
	defer conn.Close()

	bind(conn, "cn=expired,ou=users,o=example", "expired")
	passwordModify(conn, "cn=expired,ou=users,o=example", "expired", "abcd1234")


}
