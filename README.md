[![Build Status](https://travis-ci.org/MihaiLupoiu/GoDaddyDynamicDNSUpdater.svg?branch=master)](https://travis-ci.org/MihaiLupoiu/GoDaddyDynamicDNSUpdater)

# GoDaddyDynamicDNSUpdater
Script used to check and update your GoDaddy DNS server to the IP address of your current internet connection.

This client was created based on the project "GoDaddy_Powershell_DDNS" by markafox: https://github.com/markafox/GoDaddy_Powershell_DDNS

## How to use
To build the binary:
```bash 
go build .
```
To build the dockerfile:
```bash 
docker build -t godaddy-dns-updater .
```

To run the docker image:
```bash 
docker run -v $(PWD)/config.json:/config.json godaddy-dns-updater
```
