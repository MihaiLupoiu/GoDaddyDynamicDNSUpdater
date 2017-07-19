package main

/*
First go to GoDaddy developer site to create a developer account and get your key and secret

https://developer.godaddy.com/getstarted

Update the first 4 varriables with your information

*/

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
{
    "Domain":   "$your.domain.to.update  # your domain",
    "Name":     "$name_of_host #name of the A record to update",
    "Key":      "$key #key for godaddy developer API",
    "Secret":   "$secret #Secret for godday developer API"
}
*/

type Configuration struct {
	Domain string
	Name   string
	Key    string
	Secret string
}

// Log posible errors
var l = log.New(os.Stdout, ("[" + os.Args[0][2:] + "]: "), log.Ldate|log.Lshortfile)

func getConfigurationFile(configFile string) Configuration {
	configuration := Configuration{}
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		l.Println("No config.json file to read!")
	} else {
		file, _ := os.Open(configFile)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&configuration)
		if err != nil {
			l.Println("error:", err)
		}
	}
	return configuration
}

// Get public ip address there are several websites that can do this.
func getPublicIP() string {
	m := map[string]string{}
	response, err := http.Get("http://ipinfo.io/json")
	if err != nil {
		l.Println("error:", err)
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			l.Println("error:", err.Error())
		}
		err = json.Unmarshal([]byte(body), &m)
		if err != nil {
			l.Println("error:", err)
		}
	}
	return m["ip"]
}

func main() {

	configFile := flag.String("configFile", "./config.json", "JSON config file to read.")
	tmpFile := flag.String("tmpFile", "/tmp/actualIP.txt", "Path to store the last public IP.")
	test := flag.Bool("test", false, "Test URL Godaddy.")
	flag.Parse()
	fmt.Println("configFile:", *configFile)
	fmt.Println("tmpFile:", *tmpFile)
	fmt.Println("test:", *test)

	config := getConfigurationFile(*configFile)
	fmt.Println(config)

	publicIP := getPublicIP()
	fmt.Println(publicIP)

	// Check if tmp file exist, if not create.
	// Check if  IP == IP in file?
	// IF true : exit
	// Else update.

	//GODADDY Implementation to update

	client := &http.Client{}

	var url string
	if *test == true {
		url = "https://api.ote-godaddy.com/v1/domains/"
	} else {
		url = "https://api.godaddy.com/v1/domains/"
	}
	fmt.Println("URL > " + url)
	//req, _ := http.NewRequest("GET", url+config.Domain+"/records/A/"+config.Name, nil)

	//See all domains:
	//req, _ := http.NewRequest("GET", url, nil)

	// Details one domain:
	//req, _ := http.NewRequest("GET", url+config.Domain, nil)

	// GET records of one domain:
	//req, _ := http.NewRequest("GET", url+config.Domain+"/records", nil)

	// POST to recoed @ (all connections)
	// URL: = "https://api.ote-godaddy.com/v1/domains/abchub.org/records/A/%40"
	// Data to sed: [{"data": publicIP,"ttl": 600}]

	bodyToSend := map[string]interface{}{"data": publicIP, "ttl": 600}
	jsonBody, _ := json.Marshal(bodyToSend)
	req, _ := http.NewRequest("PUT", url+config.Domain+"/records/A/%40", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	//=====================
	// GLOBAL
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "sso-key "+config.Key+":"+config.Secret)

	fmt.Println("URL: ", req.URL)
	fmt.Println("Header: ", req.Header)
	fmt.Println("Body: ", req.Body)

	res, err := client.Do(req)

	var f interface{}
	if err != nil {
		l.Println("error:", err)
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.Println("error:", err.Error())
		}
		err = json.Unmarshal([]byte(body), &f)
		if err != nil {
			l.Println("error:", err)
		}
	}
	fmt.Println(f)
}
