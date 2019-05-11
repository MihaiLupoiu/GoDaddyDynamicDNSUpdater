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

type godaddyData struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int64  `json:"ttl"`
	Type string `json:"type"`
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
		log.Panic("error:", err)
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&configuration)
		if err != nil {
			log.Panic("error:", err)
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
func updateGodaddyIP(publicIP string, config Configuration) {
	bodyToSend := []map[string]interface{}{}
	data := map[string]interface{}{"data": publicIP, "ttl": 600}
	bodyToSend = append(bodyToSend, data)

	jsonBody, err := json.Marshal(bodyToSend)
	if err != nil {
		log.Println("error:", err)
	}

	if config.Name == "@" {
		config.Name = "%40"
	}
	req, err := http.NewRequest("PUT", config.URL+config.Domain+"/records/A/"+config.Name, bytes.NewBuffer(jsonBody))
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

func getGodaddyIP(config Configuration) string {
	ret := ""
	if config.Name == "@" {
		config.Name = "%40"
	}
	req, err := http.NewRequest("GET", config.URL+config.Domain+"/records/A/"+config.Name, nil)
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

		var data []godaddyData
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("error:", err)
		}
		ret = data[0].Data
	}

	return ret
}

func main() {
	// Get arguments
	configFilePath := flag.String("configFile", "./config.json", "JSON config file to read.")
	debug := flag.Bool("debug", false, "Debug mode.")
	forceUpdate := flag.Bool("force", false, "Force update.")
	flag.Parse()

	config := getConfigurationFile(*configFilePath)
	initLog("["+os.Args[0][2:]+"]: ", *debug)
	publicIP := getPublicIP()
	actualIP := getGodaddyIP(config)

	if strings.Compare(actualIP, publicIP) != 0 || *forceUpdate {
		log.Println("Changing ip of '", config.Name, ".", config.Domain, "' from:", actualIP, "to:", publicIP)
		updateGodaddyIP(publicIP, config)
	} else {
		log.Println("No IP change.")
	}

}
