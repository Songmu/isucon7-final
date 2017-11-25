#!/bin/sh

error_log() {
        MESSAGE=$1
        SCRIPT_NAME=$0
        echo "fail ${MESSAGE}" >> ./${SCRIPT_NAME}_error.log
}

setup_default_tools() {
        sudo apt-get update -y
        sudo apt-get install -y dstat sysstat wget curl unzip linux-tools strace htop ctop iftop psmisc glances || error_log "install default tools"

        wget https://github.com/tkuchiki/alp/releases/download/v0.3.1/alp_linux_amd64.zip -O alp_linux.zip \
        && sudo unzip alp_linux.zip || error_log "download alp"
        sudo install ./alp /usr/local/bin || error_log "install alp"
}

setup_percona_toolkit() {
        wget https://repo.percona.com/apt/percona-release_0.1-4.$(lsb_release -sc)_all.deb \
        && sudo dpkg -i percona-release_0.1-4.$(lsb_release -sc)_all.deb || error_log "download percona-deb-packege"
        sudo apt-get update -y
        sudo apt-get install -y percona-toolkit  || error_log "install percona-toolkit"
}

setup_default_tools
setup_percona_toolkit
