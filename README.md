# trojan
![](https://img.shields.io/github/v/release/alireza0/trojan-english.svg) 
![](https://img.shields.io/docker/pulls/alireza7/trojan-english.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alireza0/trojan-english)](https://goreportcard.com/report/github.com/alireza0/trojan-english)
[![Downloads](https://img.shields.io/github/downloads/alireza0/trojan-english/total.svg)](https://img.shields.io/github/downloads/alireza0/trojan-english/total.svg)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)


trojan multi-user management application

## Features
- Online web page and command line to manage trojan multi-users
- Start/stop/restart trojan server
- Support flow statistics and flow limit
- Command line mode management, support command completion
- Integrate acme.sh certificate application
- Generate client configuration files
- View trojan logs online in real time
- Switch between online trojan and trojan-go at any time
- Support trojan:// sharing link and QR code sharing (only for web pages)
- Support conversion to clash subscription address and import to [CLASH_FOR_WINDOWS](https://github.com/fndroid/clash_FOR_WINDOWS_PKG/releases) (only web pages)
- Limit User Period

## installation method
*For running the installation, please prepare the domain name of the server in advance*

###  a. One-click script installation

#### Installation/update
```
source <(curl -sL https://raw.githubusercontent.com/alireza0/trojan-english/master/install.sh)
```
#### Uninstall
```
source <(curl -sL https://raw.githubusercontent.com/alireza0/trojan-english/master/install.sh) --remove
```
Enter the 'trojan' to enter the management program after installation.
Browser access https://Domain_name can be available online web page to manage trojan user
Front page source code address: [train-web](https://github.com/alireza0/trojan-web)

### b. docker run
#### 1. Install mysql  

Because MariaDB memory is at least half lower than mysql, it is recommended to use the MariaDB database
```
docker run --name trojan-mariadb --restart=always -p 127.0.0.1:3306:3306 -v /home/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=trojan -e MYSQL_ROOT_HOST=% -e MYSQL_DATABASE=trojan -d mariadb:10.2
```
Both ports and root passwords and persistence directory can be changed.

#### 2. Install trojan
```
docker run -it -d --name trojan --net=host --restart=always --privileged alireza0/trojan-english
```
After running, enter the container `docker exec -it trojan bash`, Then enter 'trjan' to initialize the installation

Update management program: `source <(curl -sL https://raw.githubusercontent.com/alireza0/trojan-web/master/install.sh)`

## Run a screenshot
![avatar](asset/1.png)
![avatar](asset/2.png)

## Command Line
```
Usage:
  trojan [flags]
  trojan [command]

Available Commands:
  add           Add user
  clean         Clear designated user traffic
  completion    Automatically command to make up (support BASH and ZSH)
  del           Delete user
  help          Help about any command
  info          User information list
  log           View trojan log
  port          Modify the trojan port
  restart       Restart trojan
  start         Start up trojan
  status        View trojan status
  stop          Stop trojan
  tls           Certificate installation
  update        Update trojan
  updateWeb     Update trojan management GUI
  version       Display version number
  import [path] Import SQL file
  export [path] Export sql file
  web           Start up with web 

Flags:
  -h, --help   help for trojan
```

## Notice
After installing trojan, it is strongly recommended to open BBR and other acceleration: [Linux-NetSpeed](https://github.com/chiakge/Linux-NetSpeed)  

## Thanks
[Jrohy](https://github.com/Jrohy)