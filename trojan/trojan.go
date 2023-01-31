package trojan

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
)

// ControllMenu Trojan control menu
func ControllMenu() {
	fmt.Println()
	tType := Type()
	if tType == "trojan" {
		tType = "trojan-go"
	} else {
		tType = "trojan"
	}
	menu := []string{"Start up trojan", "Stop trojan", "Restart trojan", "View trojan status", "View Trojan log", "Modify the trojan port"}
	menu = append(menu, "switch to"+tType)
	switch util.LoopInput("Please select: ", menu, true) {
	case 1:
		Start()
	case 2:
		Stop()
	case 3:
		Restart()
	case 4:
		Status(true)
	case 5:
		go util.Log("trojan", 300)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		//block
		<-c
	case 6:
		ChangePort()
	case 7:
		if err := SwitchType(tType); err != nil {
			fmt.Println(err)
		}
	}
}

// Restart Restart trojan
func Restart() {
	util.OpenPort(core.GetConfig().LocalPort)
	util.SystemctlRestart("trojan")
}

// Start Start up trojan
func Start() {
	util.OpenPort(core.GetConfig().LocalPort)
	util.SystemctlStart("trojan")
}

// Stop Stop trojan
func Stop() {
	util.SystemctlStop("trojan")
}

// Status Get Trojan status
func Status(isPrint bool) string {
	result := util.SystemctlStatus("trojan")
	if isPrint {
		fmt.Println(result)
	}
	return result
}

// UpTime Trojan running time
func UpTime() string {
	result := strings.TrimSpace(util.ExecCommandWithResult("ps -Ao etime,args|grep -v grep|grep /usr/local/etc/trojan/config.json"))
	resultSlice := strings.Split(result, " ")
	if len(resultSlice) > 0 {
		return resultSlice[0]
	}
	return ""
}

// ChangePort Modify the trojan port
func ChangePort() {
	config := core.GetConfig()
	oldPort := config.LocalPort
	randomPort := util.RandomPort()
	fmt.Println("Current Trojan port: " + util.Green(strconv.Itoa(oldPort)))
	newPortStr := util.Input(fmt.Sprintf("Please enter the new Trojan port (if you want to use the random port %s, press Enter): ", util.Blue(strconv.Itoa(randomPort))), strconv.Itoa(randomPort))
	newPort, err := strconv.Atoi(newPortStr)
	if err != nil {
		fmt.Println("Failed to modify the port: " + err.Error())
		return
	}
	if core.WritePort(newPort) {
		util.OpenPort(newPort)
		fmt.Println(util.Green("Port modified successfully!"))
		Restart()
	} else {
		fmt.Println(util.Red("Port modified successfully!"))
	}
}

// Version Trojan version
func Version() string {
	flag := "-v"
	if Type() == "trojan-go" {
		flag = "-version"
	}
	result := strings.TrimSpace(util.ExecCommandWithResult("/usr/bin/trojan/trojan " + flag))
	if len(result) == 0 {
		return ""
	}
	firstLine := strings.Split(result, "\n")[0]
	tempSlice := strings.Split(firstLine, " ")
	return tempSlice[len(tempSlice)-1]
}

// SwitchType Switch Trojan type
func SwitchType(tType string) error {
	ARCH := runtime.GOARCH
	if ARCH != "amd64" && ARCH != "arm64" {
		return errors.New("not support " + ARCH + " machine")
	}
	if tType == "trojan" && ARCH != "amd64" {
		return errors.New("trojan not support " + ARCH + " machine")
	}
	if err := core.SetValue("trojanType", tType); err != nil {
		return err
	}
	InstallTrojan("")
	return nil
}

// Type Trojan type
func Type() string {
	tType, _ := core.GetValue("trojanType")
	if tType == "" {
		if strings.Contains(Status(false), "trojan-go") {
			tType = "trojan-go"
		} else {
			tType = "trojan"
		}
		_ = core.SetValue("trojanType", tType)
	}
	return tType
}
