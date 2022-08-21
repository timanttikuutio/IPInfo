package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"io"
	"net/http"
	"os"
	"strings"
)

type IPInfoIO struct {
	IP           string `json:"ip"`
	Hostname     string `json:"hostname"`
	Anycast      bool   `json:"anycast,omitempty"`
	Bogon        bool   `json:"bogon,omitempty"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`
	Location     string `json:"loc"`
	Organisation string `json:"org"`
	Postal       string `json:"postal"`
	Timezone     string `json:"timezone"`
}

func main() {
	args := os.Args

	prefix := color.Cyan + "[" + color.Green + "IP Info" + color.Cyan + "]" + color.Reset

	if len(args) < 2 {
		fmt.Println(prefix, color.Yellow, "No IP address provided.", color.Reset, "\n\nCommand syntax:", args[0], "<IP Address here>")
		return
	}

	req, err := http.NewRequest("GET", "https://ipinfo.io/"+args[1], nil)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	req.Header = http.Header{
		"Accept":        {"application/json"},
		"Authorization": {"Bearer 922cf6865cfa90"},
	}

	resBody, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	defer resBody.Body.Close()
	body, err := io.ReadAll(resBody.Body)

	var data IPInfoIO

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	var bogon string
	var anycast string

	if data.Bogon == true {
		bogon = "Yes"
	} else {
		bogon = "No"
	}

	if data.Anycast == true {
		anycast = "Yes"
	} else {
		anycast = "No"
	}

	fmt.Println(prefix, color.Gray, "IP:", color.Bold, data.IP, color.Reset)

	if data.Bogon == false {
		fmt.Println(prefix, color.Gray, "Hostname:", color.Bold, data.Hostname, color.Reset)
		fmt.Println(prefix, color.Gray, "Location:", color.Bold, data.City+",", data.Region+",", data.Country, color.Reset)
		fmt.Println(prefix, color.Gray, "Coordinates:", color.Bold, data.Location, color.Reset)
		fmt.Println(prefix, color.Gray, "Organisation:", color.Bold, strings.Join(strings.Split(data.Organisation, " ")[1:], " "), color.Reset)
		fmt.Println(prefix, color.Gray, "AS Number:", color.Bold, strings.Split(data.Organisation, " ")[0], color.Reset)
	}
	fmt.Println(prefix, color.Gray, "Is Bogon:", color.Bold, bogon, color.Reset)
	fmt.Println(prefix, color.Gray, "Is Anycasted:", color.Bold, anycast, color.Reset)
}
