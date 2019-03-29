# Mail URL clicker
go-guerrilla processor to point a browser to the first URL in an email

## Overview
I wrote an article https://www.linkedin.com/pulse/how-i-proved-cancel-brexit-petition-could-have-been-signed-bramsden/ where I proved that the UK petition could be signed by bots. This is the code that I used to prove that it is possible.

## Requirements
1) Server to host Go-Guerrilla
2) Domain name which has not already had a MX record set in the DNS
3) Headless Chrome 
4) Golang development language.

## Building go-guerrilla
1) Download go-guerrilla `go get github.com/flashmob/go-guerrilla`
2) Copy p_url.go into the backends directory on go-guerrilla. Not sure if this is the correct way of doing it but it worked for me.
3) In the backends directory perform a go get to get all the libraries required for p_url.go
4) I editted the Makefile in the go-guerrilla root directory as I need it create a binary for Linux, as that is what the server is running, and I added flags which removed debug information from the compiled binary to make it smaller in size. The following lines were changed:
   `GO_VARS ?= GOOS=linux`
   `LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -s -w`
5) Then I ran `make guerrillad` which produced the compiled binary file.
6) The guerrillad binary was then SFTP to the server.
7) Additional step, not needed, but as I like binarys to be small I installed UPX: `apt-get install upx`  and then ran `upx --brute guerrillad`

## Server configuration
1) A https://www.scaleway.com/ Start1-XS cloud server with Ubuntu was used which only costs €1.99 per month
2) Install chrome: 
   ```apt-get install -y unzip xvfb libxi6 libgconf-2-4
   curl -sS -o - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add
   echo "deb [arch=amd64]  http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list
   apt-get -y update
   apt-get -y install google-chrome-stable
   cd /usr/bin
   ln -s /etc/alternatives/google-chrome chrome
   ```
3) Create a linux user called chrome. This is required as Chrome does not run in headless on the Root user
4) To start Chrome on reboot add the following crontab for the chrome user: `@reboot chrome --headless --disable-gpu --remote-debugging-port=9222 https://www.chromestatus.com
5) Go-guerrilla needs a TLS public and private key to work, so as the root user run: 
`openssl req -newkey rsa:4096 -nodes -sha512 -x509 -days 3650 -nodes -out /etc/ssl/certs/mailserver.pem -keyout /etc/ssl/private/mailserver.pem`
6) Create a goguerrilla.conf.json file where guerrillad binary is located.
```{
    "log_file" : "stderr",
    "log_level" : "info",
    "allowed_hosts": [
      "bettybot.co.uk"
    ],
    "pid_file" : "/var/run/go-guerrilla.pid",
    "backend_config": {
        "log_received_mails": true,
        "save_workers_size": 1,
        "save_process" : "HeadersParser|Header|URLParser",
        "primary_mail_host" : "mail.example.com",
        "gw_save_timeout" : "30s",
        "gw_val_rcpt_timeout" : "3s"
    },
    "servers" : [
        {
            "is_enabled" : true,
            "host_name":"mail.bettybot.co.uk",
            "max_size": 1000000,
            "timeout":180,
            "listen_interface":"10.15.119.145:25",
            "max_clients": 1000,
            "log_file" : "stderr",
            "tls" : {
                "start_tls_on":true,
                "tls_always_on":false,
                "public_key_file":"/etc/ssl/certs/mailserver.pem",
                "private_key_file":"/etc/ssl/private/mailserver.pem",
                "protocols" : ["ssl3.0", "tls1.2"],
                "ciphers" : ["TLS_FALLBACK_SCSV", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305", "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305", "TLS_RSA_WITH_RC4_128_SHA", "TLS_RSA_WITH_AES_128_GCM_SHA256", "TLS_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA", "TLS_ECDHE_RSA_WITH_RC4_128_SHA", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"],
                "curves" : ["P256", "P384", "P521", "X25519"],
                "client_auth_type" : "NoClientCert"
            }
        },
        {
            "is_enabled" : false,
            "host_name":"mail.bettybot.co.uk",
            "max_size":1000000,
            "timeout":180,
            "listen_interface":"10.15.119.145:465",
            "max_clients":500,
            "log_file" : "stderr",
            "tls" : {
                "public_key_file":"/etc/ssl/certs/mailserver.pem",
                "private_key_file":"/etc/ssl/private/mailserver.pem",
                 "start_tls_on":false,
                 "tls_always_on":true,
                 "protocols" : ["ssl3.0", "tls1.2"],
                 "ciphers" : ["TLS_FALLBACK_SCSV", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305", "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305", "TLS_RSA_WITH_RC4_128_SHA", "TLS_RSA_WITH_AES_128_GCM_SHA256", "TLS_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA", "TLS_ECDHE_RSA_WITH_RC4_128_SHA", "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"],
                 "curves" : ["P256", "P384", "P521", "X25519"],
                 "client_auth_type" : "NoClientCert"
            }
        }
    ]
}
```

7) To run: `./guerrillad serve`

## Configuring DNS
For this you will need a Domain name which is currently not being used for Mail already. 
1) Create a A record with the hostname of mail and the content being the IP address of the server
2) Create a MX record with the hostname of @ and the content being mail.YOURDOMAINNAME. E.G. in my example it is mail.bettybot.co.uk

## To test
To test just send an email to any email address for you domain with a URL in the body of the message. When it is received by Go-Guerrillad you should see a lot of information from chrome and on the last line you will see the HTML page title of the URL.

