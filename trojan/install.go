package trojan

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
	"trojan/asset"
	"trojan/core"
	"trojan/util"
)

var (
	dockerInstallUrl = "https://docker-install.netlify.app/install.sh"
	dbDockerRun      = "docker run --name trojan-mariadb --restart=always -p %d:3306 -v /home/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=%s -e MYSQL_ROOT_HOST=%% -e MYSQL_DATABASE=trojan -d mariadb:10.2"
)

// InstallMenu installation manual
func InstallMenu() {
	fmt.Println()
	menu := []string{"Update trojan", "Certificate application", "Install mysql"}
	switch util.LoopInput("Please select: ", menu, true) {
	case 1:
		InstallTrojan("")
	case 2:
		InstallTls()
	case 3:
		InstallMysql()
	default:
		return
	}
}

// InstallDocker Install docker
func InstallDocker() {
	if !util.CheckCommandExists("docker") {
		util.RunWebShell(dockerInstallUrl)
		fmt.Println()
	}
}

// InstallTrojan Install trojan
func InstallTrojan(version string) {
	fmt.Println()
	data := string(asset.GetAsset("trojan-install.sh"))
	checkTrojan := util.ExecCommandWithResult("systemctl list-unit-files|grep trojan.service")
	if (checkTrojan == "" && runtime.GOARCH != "amd64") || Type() == "trojan-go" {
		data = strings.ReplaceAll(data, "TYPE=0", "TYPE=1")
	}
	if version != "" {
		data = strings.ReplaceAll(data, "INSTALL_VERSION=\"\"", "INSTALL_VERSION=\""+version+"\"")
	}
	util.ExecCommand(data)
	util.OpenPort(443)
	util.SystemctlRestart("trojan")
	util.SystemctlEnable("trojan")
}

// InstallTls Install certificate
func InstallTls() {
	domain := ""
	server := "letsencrypt"
	fmt.Println()
	choice := util.LoopInput("Please choose the method of using the certificate: ", []string{"Let's Encrypt Certificate", "ZeroSSL Certificate", "BuyPass Certificate", "Custom certificate path"}, true)
	if choice < 0 {
		return
	} else if choice == 4 {
		crtFile := util.Input("Please enter the Cert File path of Certificate: ", "")
		keyFile := util.Input("Please enter the key file path of Certificate: ", "")
		if !util.IsExists(crtFile) || !util.IsExists(keyFile) {
			fmt.Println("The input CERT or Key File does not exist!")
		} else {
			domain = util.Input("Please enter the domain name corresponding to this Certification: ", "")
			if domain == "" {
				fmt.Println("Enter domain name empty!")
				return
			}
			core.WriteTls(crtFile, keyFile, domain)
		}
	} else {
		if choice == 2 {
			server = "zerossl"
		} else if choice == 3 {
			server = "buypass"
		}
		localIP := util.GetLocalIP()
		fmt.Printf("This machine IP: %s\n", localIP)
		for {
			domain = util.Input("Please enter the domain name to apply for Certification: ", "")
			ipList, err := net.LookupIP(domain)
			fmt.Printf("%s Analysis IP: %v\n", domain, ipList)
			if err != nil {
				fmt.Println(err)
				fmt.Println("The domain name is wrong, please re-enter")
				continue
			}
			checkIp := false
			for _, ip := range ipList {
				if localIP == ip.String() {
					checkIp = true
				}
			}
			if checkIp {
				break
			} else {
				fmt.Println("The input domain name is inconsistent with the IP of this machine, please re-enter!")
			}
		}
		util.InstallPack("socat")
		if !util.IsExists("/root/.acme.sh/acme.sh") {
			util.RunWebShell("https://get.acme.sh")
		}
		util.SystemctlStop("trojan-web")
		util.OpenPort(80)
		checkResult := util.ExecCommandWithResult("/root/.acme.sh/acme.sh -v|tr -cd '[0-9]'")
		acmeVersion, _ := strconv.Atoi(checkResult)
		if acmeVersion < 300 {
			util.ExecCommand("/root/.acme.sh/acme.sh --upgrade")
		}
		if server != "letsencrypt" {
			var email string
			for {
				email = util.Input(fmt.Sprintf("Please enter the mailbox required to apply for a%S domain name: ", server), "")
				if email == "" {
					fmt.Println("The mailbox address of the domain name is empty!")
					return
				} else if util.VerifyEmailFormat(email) {
					break
				} else {
					fmt.Println("The mailbox format is incorrect, please re-enter!")
				}
			}
			util.ExecCommand(fmt.Sprintf("bash /root/.acme.sh/acme.sh --server %s --register-account -m %s", server, email))
		}
		issueCommand := fmt.Sprintf("bash /root/.acme.sh/acme.sh --issue -d %s --debug --standalone --keylength ec-256 --force --server %s", domain, server)
		if server == "buypass" {
			issueCommand = issueCommand + " --days 170"
		}
		util.ExecCommand(issueCommand)
		crtFile := "/root/.acme.sh/" + domain + "_ecc" + "/fullchain.cer"
		keyFile := "/root/.acme.sh/" + domain + "_ecc" + "/" + domain + ".key"
		core.WriteTls(crtFile, keyFile, domain)
	}
	Restart()
	util.SystemctlRestart("trojan-web")
	fmt.Println()
}

// InstallMysql Install mysql
func InstallMysql() {
	var (
		mysql  core.Mysql
		choice int
	)
	fmt.Println()
	if util.IsExists("/.dockerenv") {
		choice = 2
	} else {
		choice = util.LoopInput("Please select: ", []string{"Install the docker version of mysql(mariadb)", "Enter custom MYSQL connection"}, true)
	}
	if choice < 0 {
		return
	} else if choice == 1 {
		mysql = core.Mysql{ServerAddr: "127.0.0.1", ServerPort: util.RandomPort(), Password: util.RandString(5), Username: "root", Database: "trojan"}
		InstallDocker()
		fmt.Println(fmt.Sprintf(dbDockerRun, mysql.ServerPort, mysql.Password))
		if util.CheckCommandExists("setenforce") {
			util.ExecCommand("setenforce 0")
		}
		util.OpenPort(mysql.ServerPort)
		util.ExecCommand(fmt.Sprintf(dbDockerRun, mysql.ServerPort, mysql.Password))
		db := mysql.GetDB()
		for {
			fmt.Printf("%s mariadb Start up is in progress, please wait a little...\n", time.Now().Format("2006-01-02 15:04:05"))
			err := db.Ping()
			if err == nil {
				db.Close()
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
		fmt.Println("mariadb is now up and running!")
	} else if choice == 2 {
		mysql = core.Mysql{}
		for {
			for {
				mysqlUrl := util.Input("Please enter the MySQL connection address (format: host:port). Press Enter to use default address (127.0.0.1:3306), otherwise enter a custom connection address: ",
					"127.0.0.1:3306")
				urlInfo := strings.Split(mysqlUrl, ":")
				if len(urlInfo) != 2 {
					fmt.Printf("The input %s does not match the matching format (host:port) \n", mysqlUrl)
					continue
				}
				port, err := strconv.Atoi(urlInfo[1])
				if err != nil {
					fmt.Printf("%s is not a number\n", urlInfo[1])
					continue
				}
				mysql.ServerAddr, mysql.ServerPort = urlInfo[0], port
				break
			}
			mysql.Username = util.Input("Please enter the username of MySQL (Press Enter for root): ", "root")
			mysql.Password = util.Input(fmt.Sprintf("Please enter the password of mysql user %s: ", mysql.Username), "")
			db := mysql.GetDB()
			if db != nil && db.Ping() == nil {
				mysql.Database = util.Input("Please enter the database name used (there is no existence can be created automatically. Press Enter to use 'trojan' as default.): ", "trojan")
				db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", mysql.Database))
				break
			} else {
				fmt.Println("Failed to connect mysql, please re-enter")
			}
		}
	}
	mysql.CreateTable()
	core.WriteMysql(&mysql)
	if userList, _ := mysql.GetData(); len(userList) == 0 {
		AddUser()
	}
	Restart()
	fmt.Println()
}
