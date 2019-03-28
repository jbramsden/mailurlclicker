# Mail URL clicker
go-guerrilla processor to point a browser to the first URL in an email

## Overview
I wrote an article https://www.linkedin.com/pulse/how-i-proved-cancel-brexit-petition-could-have-been-signed-bramsden/ where I proved the the UK petition could be signed by bots. This is the code that I used to prove that it is possible.

## Requirements
1) Server to host Go-Guerrilla
2) Domain name which has not already had a MX record set in the DNS
3) Headless Chrome 

## Building go-guerrilla
1) Download go-guerrilla 
2) Copy p_url.go into the backends directory on go-guerrilla. Not sure if this is the correct way of doing it but it worked for me.
3) In the backends directory perform a go get to get all the libraries required p_url.go
4) I editted the Makefile in the go-guerrilla root directory as I need it create a binary for Linux, as that is what the server is running, and I added flags which removed debug information from the compiled binary to make it smaller in size. The following lines were changed:
   GO_VARS ?= GOOS=linux
   LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -s -w
5) Then I ran make guerrillad which produced the compiled binary file.
6) The guerrillad binary was then SFTP to the server.
7) Additional step, not needed, but as I like binarys to be small I installed UPX: apt-get install upx  and then ran upx --brute guerrillad

## Server configuration
1) A https://www.scaleway.com/ Start1-XS cloud server with Ubuntu was used which only costs €1.99 per month
2) Install chrome: 
   apt-get install -y unzip xvfb libxi6 libgconf-2-4
   curl -sS -o - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add
   echo "deb [arch=amd64]  http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list
   apt-get -y update
   apt-get -y install google-chrome-stable
   cd /usr/bin
   ln -s /etc/alternatives/google-chrome chrome
3) Create a linux user called chrome. This is required as Chrome does not run in headless on the Root user
4) To start Chrome on reboot add the following crontab for the chrome user: @reboot chrome --headless --disable-gpu --remote-debugging-port=9222 https://www.chromestatus.com
5) 

