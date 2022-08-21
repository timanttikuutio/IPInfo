package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/joho/godotenv"
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

	configPath, _ := os.UserConfigDir()

	err := godotenv.Load(configPath + "/ipinfoConfig.env")
	if err != nil {
		fmt.Println(prefix, color.Yellow, "The config file does not exist. Generating one at:", configPath+"/ipinfoConfig.env", color.Reset)

		file, err := os.Create(configPath + "/ipinfoConfig.env")
		if err != nil {
			fmt.Println(prefix, color.Reset, "An error occurred whilst creating the config file.", color.Reset)
		}

		_, err = file.WriteString("IPINFO_APIKEY=")
		if err != nil {
			return
		}
	}

	if len(args) < 2 {
		fmt.Println(prefix, color.Yellow, "No IP address provided.", color.Reset, "\n\nCommand syntax:", args[0],
			"<IP Address here>")
		return
	}

	req, err := http.NewRequest("GET", "https://ipinfo.io/"+args[1], nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header = http.Header{"Accept": {"application/json"}}

	apikey, exists := os.LookupEnv("IPINFO_APIKEY")
	if !exists || apikey == "" {
		fmt.Println(prefix, color.Yellow, "The API-key seems to be unset, accessing ipinfo.io as a guest.\n",
			strings.Repeat(" ", 9), "NOTE: heavy rate-limits will apply. Consider creating a free API-key at https://ipinfo.io")
	} else {
		req.Header.Add("Authorization", "Bearer "+apikey)
	}

	resBody, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resBody.StatusCode == 403 {
		fmt.Println(prefix, color.Red, "Your ipinfo.io API key seems to be invalid. Please check the key or omit "+
			"it, to run the script as guest.\n", strings.Repeat(" ", 9), "NOTE: You may be rate-limited "+
			"rather quickly without a valid API key.")
		return
	}
	if resBody.StatusCode == 429 {
		fmt.Println(
			prefix, color.Yellow,
			"It seems like you've exceeded the rate-limit. Please try again at later date.", color.Reset)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resBody.Body)
	body, err := io.ReadAll(resBody.Body)

	var data IPInfoIO

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
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
