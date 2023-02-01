package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// PortIsUse Determine whether the port is occupied
func PortIsUse(port int) bool {
	_, tcpError := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Millisecond*50)
	udpAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	udpConn, udpError := net.ListenUDP("udp", udpAddr)
	if udpConn != nil {
		defer udpConn.Close()
	}
	return tcpError == nil || udpError != nil
}

// RandomPort Get the random port that is not occupied
func RandomPort() int {
	for {
		rand.Seed(time.Now().UnixNano())
		newPort := rand.Intn(65536)
		if !PortIsUse(newPort) {
			return newPort
		}
	}
}

// IsExists Check whether the specified path file or file clip exists
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat Get file information
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GetLocalIP Get the IPv4 address of this machine
func GetLocalIP() string {
	resp, err := http.Get("http://api.ipify.org")
	if err != nil {
		resp, _ = http.Get("http://icanhazip.com")
	}
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	return string(s)
}

// InstallPack Install the specified name software
func InstallPack(name string) {
	if !CheckCommandExists(name) {
		if CheckCommandExists("yum") {
			ExecCommand("yum install -y " + name)
		} else if CheckCommandExists("apt-get") {
			ExecCommand("apt-get update")
			ExecCommand("apt-get install -y " + name)
		}
	}
}

// OpenPort Open the specified port
func OpenPort(port int) {
	if CheckCommandExists("firewall-cmd") {
		ExecCommand(fmt.Sprintf("firewall-cmd --zone=public --add-port=%d/tcp --add-port=%d/udp --permanent >/dev/null 2>&1", port, port))
		ExecCommand("firewall-cmd --reload >/dev/null 2>&1")
	} else {
		if len(ExecCommandWithResult(fmt.Sprintf(`iptables -nvL --line-number|grep -w "%d"`, port))) > 0 {
			return
		}
		ExecCommand(fmt.Sprintf("iptables -I INPUT -p tcp --dport %d -j ACCEPT", port))
		ExecCommand(fmt.Sprintf("iptables -I INPUT -p udp --dport %d -j ACCEPT", port))
		ExecCommand(fmt.Sprintf("iptables -I OUTPUT -p udp --sport %d -j ACCEPT", port))
		ExecCommand(fmt.Sprintf("iptables -I OUTPUT -p tcp --sport %d -j ACCEPT", port))
	}
}

// Log Real-time printing specified service log
func Log(serviceName string, line int) {
	result, _ := LogChan(serviceName, "-n "+strconv.Itoa(line), make(chan byte))
	for line := range result {
		fmt.Println(line)
	}
}

// LogChan Specify the real-time log, return to chan
func LogChan(serviceName, param string, closeChan chan byte) (chan string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("journalctl -f -u %s -o cat %s", serviceName, param))

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err: ", err.Error())
		return nil, err
	}
	ch := make(chan string, 100)
	stdoutScan := bufio.NewScanner(stdout)
	go func() {
		for stdoutScan.Scan() {
			select {
			case <-closeChan:
				stdout.Close()
				return
			default:
				ch <- stdoutScan.Text()
			}
		}
	}()
	return ch, nil
}
