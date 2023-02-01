package trojan

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
)

// UserMenu User management menu
func UserMenu() {
	fmt.Println()
	menu := []string{"Add User", "Delete User", "Limit Traffic", "Clear Traffic", "Set Time Limit", "Cancel Time Limit"}
	switch util.LoopInput("Please select: ", menu, false) {
	case 1:
		AddUser()
	case 2:
		DelUser()
	case 3:
		SetUserQuota()
	case 4:
		CleanData()
	case 5:
		SetupExpire()
	case 6:
		CancelExpire()
	}
}

// AddUser Add user
func AddUser() {
	randomUser := util.RandString(4)
	randomPass := util.RandString(8)
	inputUser := util.Input(fmt.Sprintf("Press Enter to use random generated username: %s, otherwise enter a custom user name: ", randomUser), randomUser)
	if inputUser == "admin" {
		fmt.Println(util.Yellow("Cannot create a new user with username 'admin'!"))
		return
	}
	mysql := core.GetMysql()
	if user := mysql.GetUserByName(inputUser); user != nil {
		fmt.Println(util.Yellow("The username entered: " + inputUser + " exists!"))
		return
	}
	inputPass := util.Input(fmt.Sprintf("Press Enter to use random generated password: %s, otherwise enter the custom password: ", randomPass), randomPass)
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass))
	if user := mysql.GetUserByPass(base64Pass); user != nil {
		fmt.Println(util.Yellow("The password entered: " + inputPass + " exists!"))
		return
	}
	if mysql.CreateUser(inputUser, base64Pass, inputPass) == nil {
		fmt.Println("New user created successfully!")
	}
}

// DelUser Delete user
func DelUser() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Please select the user number to be deleted: ", userList, true)
	if choice == -1 {
		return
	}
	if mysql.DeleteUser(userList[choice-1].ID) == nil {
		fmt.Println("User deleted!")
		Restart()
	}
}

// SetUserQuota Restricted user traffic
func SetUserQuota() {
	var (
		limit int
		err   error
	)
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Please select the user number to limit the traffic: ", userList, true)
	if choice == -1 {
		return
	}
	for {
		quota := util.Input("Please enter the user"+userList[choice-1].Username+"Limited traffic (unit byte)", "")
		limit, err = strconv.Atoi(quota)
		if err != nil {
			fmt.Printf("%s is not a number, please re-enter!\n", quota)
		} else {
			break
		}
	}
	if mysql.SetQuota(userList[choice-1].ID, limit) == nil {
		fmt.Println("Restricting traffic " + userList[choice-1].Username + " is set to " + util.Bytefmt(uint64(limit)))
	}
}

// CleanData Clear user traffic
func CleanData() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Please select the user number to clear the flow: ", userList, true)
	if choice == -1 {
		return
	}
	if mysql.CleanData(userList[choice-1].ID) == nil {
		fmt.Println("Traffic cleared successfully!")
	}
}

// CancelExpire Cancel
func CancelExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Please select the user number to cancel the time limit: ", userList, true)
	if choice == -1 {
		return
	}
	if userList[choice-1].UseDays == 0 {
		fmt.Println(util.Yellow("The selected user is not set to time limit!"))
		return
	}
	if mysql.CancelExpire(userList[choice-1].ID) == nil {
		fmt.Println("Time-limit canceled successful!")
	}
}

// SetupExpire Set time
func SetupExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Please select the user number to set the time limit: ", userList, true)
	if choice == -1 {
		return
	}
	useDayStr := util.Input("Please enter the number of days to be used to limit: ", "")
	if useDayStr == "" {
		return
	} else if strings.Contains(useDayStr, "-") {
		fmt.Println(util.Yellow("The number of days cannot be negative"))
		return
	} else if !util.IsInteger(useDayStr) {
		fmt.Println(util.Yellow("Input to non-integer!"))
		return
	}
	useDays, _ := strconv.Atoi(useDayStr)
	if mysql.SetExpire(userList[choice-1].ID, uint(useDays)) == nil {
		fmt.Println("The time limit is set!")
	}
}

// CleanDataByName Clear designated user traffic
func CleanDataByName(usernames []string) {
	mysql := core.GetMysql()
	if err := mysql.CleanDataByName(usernames); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Clear traffic successfully!")
	}
}

// UserList Get the user list and print the display
func UserList(ids ...string) []*core.User {
	mysql := core.GetMysql()
	userList, err := mysql.GetData(ids...)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	domain, port := GetDomainAndPort()
	for i, k := range userList {
		pass, err := base64.StdEncoding.DecodeString(k.Password)
		if err != nil {
			pass = []byte("")
		}
		fmt.Printf("%d.\n", i+1)
		fmt.Println("username: " + k.Username)
		fmt.Println("password: " + string(pass))
		fmt.Println("Upload traffic: " + util.Cyan(util.Bytefmt(k.Upload)))
		fmt.Println("Download traffic: " + util.Cyan(util.Bytefmt(k.Download)))
		if k.Quota < 0 {
			fmt.Println("Flow limit: " + util.Cyan("Unlimited"))
		} else {
			fmt.Println("Flow limit: " + util.Cyan(util.Bytefmt(uint64(k.Quota))))
		}
		if k.UseDays == 0 {
			fmt.Println("Date of Expiry: " + util.Cyan("Unlimited"))
		} else {
			fmt.Println("Date of Expiry: " + util.Cyan(k.ExpiryDate))
		}
		remark := url.QueryEscape(k.Username)
		fmt.Println("Share link: " + util.Green(fmt.Sprintf("trojan://%s@%s:%d#%s", string(pass), domain, port, remark)))
		fmt.Println()
	}
	return userList
}
