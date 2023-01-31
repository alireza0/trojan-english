package trojan

import (
	"crypto/sha256"
	"fmt"
	"trojan/core"
	"trojan/util"
)

// WebMenu web management menu
func WebMenu() {
	fmt.Println()
	menu := []string{"Reset the web administrator password", "Modify the displayed domain name (not apply for a Certificate)"}
	switch util.LoopInput("Please choose: ", menu, true) {
	case 1:
		ResetAdminPass()
	case 2:
		SetDomain("")
	}
}

// ResetAdminPass Reset the administrator password
func ResetAdminPass() {
	inputPass := util.Input("Please enter the admin user password: ", "")
	if inputPass == "" {
		fmt.Println("Changes rejected!")
	} else {
		encryPass := sha256.Sum224([]byte(inputPass))
		err := core.SetValue("admin_pass", fmt.Sprintf("%x", encryPass))
		if err == nil {
			fmt.Println(util.Green("The admin password is not reset!"))
		} else {
			fmt.Println(err)
		}
	}
}

// SetDomain Set the displayed domain name
func SetDomain(domain string) {
	if domain == "" {
		domain = util.Input("Please enter the domain name address to be displayed: ", "")
	}
	if domain == "" {
		fmt.Println("Changes rejected!")
	} else {
		core.WriteDomain(domain)
		Restart()
		fmt.Println("Modify domain is done successfully!")
	}
}

// GetDomainAndPort Get domain name and port
func GetDomainAndPort() (string, int) {
	config := core.GetConfig()
	return config.SSl.Sni, config.LocalPort
}
