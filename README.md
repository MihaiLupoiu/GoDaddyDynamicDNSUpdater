[![Build Status](https://travis-ci.org/MihaiLupoiu/GoDaddyDynamicDNSUpdater.svg?branch=master)](https://travis-ci.org/MihaiLupoiu/GoDaddyDynamicDNSUpdater)

# GoDaddyDynamicDNSUpdater
Script used to check and update your GoDaddy DNS server to the IP address of your current internet connection.

This client was created based on the project "GoDaddy_Powershell_DDNS" by markafox: https://github.com/markafox/GoDaddy_Powershell_DDNS

## Configuration

In order to use the script it is necesary to generate a `key` and `secret` from the godaddy developer's page.

Example of configuration using the godaddy testing api:
```json
{
    "URL":"https://api.ote-godaddy.com/v1/domains/",
    "Domain":"abchub.org",
    "Name":"@",
    "Key":"UzQxLikm_46KxDFnbjN7cQjmw6wocia",
    "Secret":"46L26ydpkwMaKZV6uVdDWe"
}
```

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
test_config.json

## Kubernetes
Deploy in kubernetes:
```bash
kubectl create configmap godaddy-dynamic-dns-updater-conf --from-file=config.json
kubectl apply -f ./cronjob.yaml
```

## TODO:
* Eliminate the file and compare the IP from Godaddy with my public IP. Also show in the logs the IP before and after the change.
