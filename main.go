package main

/*
First go to GoDaddy developer site to create a developer account and get your key and secret

https://developer.godaddy.com/getstarted

Update the first 4 varriables with your information

*/
import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"
)

// Configuration struct for reading the config.json
type Configuration struct {
	URL    string
	Domain string
	Name   string
	Key    string
	Secret string
}

// init log level to print in syslog
func initLog(name string, debug bool) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if debug {
		log.SetOutput(os.Stdout)
	} else {
		var logWriter, logErr = syslog.New(syslog.LOG_ERR, string(os.Args[0][2:]))
		if logErr == nil {
			log.SetOutput(logWriter)
		}
	}
}

// Parshe json configuration file.
func getConfigurationFile(configFile string) Configuration {
	configuration := Configuration{}
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("error:", err)
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&configuration)
		if err != nil {
			log.Println("error:", err)
		}
	}
	return configuration
}

// Get public ip address there are several websites that can do this.
func getPublicIP() string {
	var strIP string
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Println("error:", err)
	} else {
		ip, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error:", err)
		} else {
			strIP = string(ip)
		}
	}
	return strIP
}

// Godaddy Implementation to update public IP
func updateIPGodaddy(url, publicIP string, config Configuration) {

	bodyToSend := []map[string]interface{}{}
	data := map[string]interface{}{"data": publicIP, "ttl": 600}
	bodyToSend = append(bodyToSend, data)

	jsonBody, _ := json.Marshal(bodyToSend)
	if config.Name == "@" {
		config.Name = "%40"
	}
	req, err := http.NewRequest("PUT", url+config.Domain+"/records/A/"+config.Name, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("error:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "sso-key "+config.Key+":"+config.Secret)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Println("error:", err)
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error:", err.Error())
		}
		strBody := string([]byte(body))
		if len(strBody) > 2 {
			log.Println(strBody)
		}
	}
}

// Get temporal ip sored in the /tmp/actualFile.txt
func getTmpIP(path string) string {
	tmpFile, err := os.Open(path)
	var tmpIP string
	if err == nil {
		readFile := bufio.NewReader(tmpFile)
		tmpIP, err = readFile.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("error: No IP in", path, ".")
			} else {
				log.Println("error:", err)
			}
		}
		if err := tmpFile.Close(); err != nil {
			log.Println("error:", err)
		}

	} else {
		log.Println("error:", err)
	}
	return strings.TrimRight(tmpIP, "\n")
}

// Update IP in the /tmp/actualFile.txt.
func updateTmpIP(path, publicIP string) {
	tmpFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("error:", err)
	}
	_, err = tmpFile.WriteString(publicIP + "\n")
	tmpFile.Sync()
	if err := tmpFile.Close(); err != nil {
		log.Println("error:", err)
	}
}

func main() {
	// Get arguments
	configFilePath := flag.String("configFile", "./config.json", "JSON config file to read.")
	tmpFilePath := flag.String("tmpFile", "/tmp/actualIP.txt", "Path to store the last public IP.")
	debug := flag.Bool("debug", false, "Debug mode.")
	flag.Parse()

	config := getConfigurationFile(*configFilePath)
	initLog("["+os.Args[0][2:]+"]: ", *debug)
	publicIP := getPublicIP()
	tmpIP := getTmpIP(*tmpFilePath)

	if strings.Compare(tmpIP, publicIP) != 0 {
		updateIPGodaddy(config.URL, publicIP, config)
		updateTmpIP(*tmpFilePath, publicIP)
		log.Println("IP updated.")
	} else {
		log.Println("No IP change.")
	}
}
