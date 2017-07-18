package main

/*
First go to GoDaddy developer site to create a developer account and get your key and secret

https://developer.godaddy.com/getstarted
 
Update the first 4 varriables with your information

*/

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
    Domain	string
	Name   	string
	Key		string
	Secret 	string
}

/* Add option to read from commands like -c /path/to/file.json. */
func main() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration.Domain)
	
	// Get public ip address there are several websites that can do this.
	/*
	{
		"ip": "90.106.193.21",
		"hostname": "21.pool90-106-193.dynamic.orange.es",
		"city": "Port Colom",
		"region": "Islas Baleares",
		"country": "ES",
		"loc": "39.4192,3.2600",
		"org": "AS12479 Orange Espagne SA"
	}
	*/
	response, err := http.Get("http://ipinfo.io/json")
	if err != nil {
			fmt.Println(err)
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err.Error())
		}
		m := map[string]string{}
		err = json.Unmarshal([]byte(body), &m)
		if err != nil {
			panic(err)
		}
		fmt.Println(m["ip"])
	}

}