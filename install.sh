#!/bin/bash
# Author: Jrohy
# github: https://github.com/Jrohy/trojan

#Define the operation variable, 0 is no, 1 is yes
HELP=0

REMOVE=0

UPDATE=0

DOWNLAOD_URL="https://github.com/Jrohy/trojan/releases/download/"

VERSION_CHECK="https://api.github.com/repos/Jrohy/trojan/releases/latest"

SERVICE_URL="https://raw.githubusercontent.com/Jrohy/trojan/master/asset/trojan-web.service"

[[ -e /var/lib/trojan-manager ]] && UPDATE=1

#Centos Temporary cancel
[[ -f /etc/redhat-release && -z $(echo $SHELL|grep zsh) ]] && unalias -a

[[ -z $(echo $SHELL|grep zsh) ]] && SHELL_WAY="bash" || SHELL_WAY="zsh"

#######color code########
RED="31m"
GREEN="32m"
YELLOW="33m"
BLUE="36m"
FUCHSIA="35m"

colorEcho(){
    COLOR=$1
    echo -e "\033[${COLOR}${@:2}\033[0m"
}

#######get params#########
while [[ $# > 0 ]];do
    KEY="$1"
    case $KEY in
        --remove)
        REMOVE=1
        ;;
        -h|--help)
        HELP=1
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done
#############################

help(){
    echo "bash $0 [-h|--help] [--remove]"
    echo "  -h, --help           Show help"
    echo "      --remove         remove trojan"
    return 0
}

removeTrojan() {
    #Remove Trojan
    rm -rf /usr/bin/trojan >/dev/null 2>&1
    rm -rf /usr/local/etc/trojan >/dev/null 2>&1
    rm -f /etc/systemd/system/trojan.service >/dev/null 2>&1

    #Remove Trojan management GUI
    rm -f /usr/local/bin/trojan >/dev/null 2>&1
    rm -rf /var/lib/trojan-manager >/dev/null 2>&1
    rm -f /etc/systemd/system/trojan-web.service >/dev/null 2>&1

    systemctl daemon-reload

    #Remove trojan's private db
    docker rm -f trojan-mysql trojan-mariadb >/dev/null 2>&1
    rm -rf /home/mysql /home/mariadb >/dev/null 2>&1
    
    #Removing environment variables
    sed -i '/trojan/d' ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc

    colorEcho ${GREEN} "uninstall success!"
}

checkSys() {
    #Check whether it is root
    [ $(id -u) != "0" ] && { colorEcho ${RED} "Error: You must be root to run this script"; exit 1; }

    ARCH=$(uname -m 2> /dev/null)
    if [[ $ARCH != x86_64 && $ARCH != aarch64 ]];then
        colorEcho $YELLOW "not support $ARCH machine".
        exit 1
    fi

    if [[ `command -v apt-get` ]];then
        PACKAGE_MANAGER='apt-get'
    elif [[ `command -v dnf` ]];then
        PACKAGE_MANAGER='dnf'
    elif [[ `command -v yum` ]];then
        PACKAGE_MANAGER='yum'
    else
        colorEcho $RED "Not support OS!"
        exit 1
    fi

    # Automatically added when the /usr/local/bin path is missing
    [[ -z `echo $PATH|grep /usr/local/bin` ]] && { echo 'export PATH=$PATH:/usr/local/bin' >> /etc/bashrc; source /etc/bashrc; }
}

#Install dependencies
installDependent(){
    if [[ ${PACKAGE_MANAGER} == 'dnf' || ${PACKAGE_MANAGER} == 'yum' ]];then
        ${PACKAGE_MANAGER} install socat crontabs bash-completion -y
    else
        ${PACKAGE_MANAGER} update
        ${PACKAGE_MANAGER} install socat cron bash-completion xz-utils -y
    fi
}

setupCron() {
    if [[ `crontab -l 2>/dev/null|grep acme` ]]; then
        if [[ -z `crontab -l 2>/dev/null|grep trojan-web` || `crontab -l 2>/dev/null|grep trojan-web|grep "&"` ]]; then
            #Calculate the actual time of VPS at 3 am in Beijing time
            ORIGIN_TIME_ZONE=$(date -R|awk '{printf"%d",$6}')
            LOCAL_TIME_ZONE=${ORIGIN_TIME_ZONE%00}
            BEIJING_ZONE=8
            BEIJING_UPDATE_TIME=3
            DIFF_ZONE=$[$BEIJING_ZONE-$LOCAL_TIME_ZONE]
            LOCAL_TIME=$[$BEIJING_UPDATE_TIME-$DIFF_ZONE]
            if [ $LOCAL_TIME -lt 0 ];then
                LOCAL_TIME=$[24+$LOCAL_TIME]
            elif [ $LOCAL_TIME -ge 24 ];then
                LOCAL_TIME=$[$LOCAL_TIME-24]
            fi
            crontab -l 2>/dev/null|sed '/acme.sh/d' > crontab.txt
            echo "0 ${LOCAL_TIME}"' * * * systemctl stop trojan-web; "/root/.acme.sh"/acme.sh --cron --home "/root/.acme.sh" > /dev/null; systemctl start trojan-web' >> crontab.txt
            crontab crontab.txt
            rm -f crontab.txt
        fi
    fi
}

installTrojan(){
    local SHOW_TIP=0
    if [[ $UPDATE == 1 ]];then
        systemctl stop trojan-web >/dev/null 2>&1
        rm -f /usr/local/bin/trojan
    fi
    LASTEST_VERSION=$(curl -H 'Cache-Control: no-cache' -s "$VERSION_CHECK" | grep 'tag_name' | cut -d\" -f4)
    echo "downloading management GUI `colorEcho $BLUE $LASTEST_VERSION` version..."
    [[ $ARCH == x86_64 ]] && BIN="trojan-linux-amd64" || BIN="trojan-linux-arm64" 
    curl -L "$DOWNLAOD_URL/$LASTEST_VERSION/$BIN" -o /usr/local/bin/trojan
    chmod +x /usr/local/bin/trojan
    if [[ ! -e /etc/systemd/system/trojan-web.service ]];then
        SHOW_TIP=1
        curl -L $SERVICE_URL -o /etc/systemd/system/trojan-web.service
        systemctl daemon-reload
        systemctl enable trojan-web
    fi
    #Command to complete environment variables
    [[ -z $(grep trojan ~/.${SHELL_WAY}rc) ]] && echo "source <(trojan completion ${SHELL_WAY})" >> ~/.${SHELL_WAY}rc
    source ~/.${SHELL_WAY}rc
    if [[ $UPDATE == 0 ]];then
        colorEcho $GREEN "Install trojan management GUI successfully!\n"
        echo -e "Run the command `colorEcho $BLUE trojan` for trojan management\n"
        /usr/local/bin/trojan
    else
        if [[ `cat /usr/local/etc/trojan/config.json|grep -w "\"db\""` ]];then
            sed -i "s/\"db\"/\"database\"/g" /usr/local/etc/trojan/config.json
            systemctl restart trojan
        fi
        /usr/local/bin/trojan upgrade db
        if [[ -z `cat /usr/local/etc/trojan/config.json|grep sni` ]];then
            /usr/local/bin/trojan upgrade config
        fi
        systemctl restart trojan-web
        colorEcho $GREEN "Update trojan management GUI succeed!\n"
    fi
    setupCron
    [[ $SHOW_TIP == 1 ]] && echo "Browser access'`colorEcho $BLUE https://domain name`' can be managed by online trojan multi-user"
}

main(){
    [[ ${HELP} == 1 ]] && help && return
    [[ ${REMOVE} == 1 ]] && removeTrojan && return
    [[ $UPDATE == 0 ]] && echo "Installing trojan management GUI..." || echo "Updating trojan management GUI..."
    checkSys
    [[ $UPDATE == 0 ]] && installDependent
    installTrojan
}

main