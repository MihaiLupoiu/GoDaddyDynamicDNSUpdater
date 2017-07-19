package main

/*
First go to GoDaddy developer site to create a developer account and get your key and secret

https://developer.godaddy.com/getstarted

Update the first 4 varriables with your information

*/

import (
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

	fmt.Println("configFile:", *configFile)
	fmt.Println("tmpFile:", *tmpFile)

	config := getConfigurationFile(*configFile)
	fmt.Println(config)

	publicIP := getPublicIP()
	fmt.Println(publicIP)

	// Check if tmp file exist, if not create.
	// Check if  IP == IP in file?
	// IF true : exit
	// Else update.

	//GODADDY Implementation to update

}
