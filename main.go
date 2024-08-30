package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type IPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

func main() {
	CheckOS()
	IP, City := GetPublic()

	fmt.Printf("--------------------------------------------------------------------\n")
	fmt.Printf("| Network Service | Local IPv4 Address |  Public IP   |    City    |\n")
	fmt.Printf("--------------------------------------------------------------------\n")

	services := GetNetworkServices()

	for _, service := range services {
		info := GetNetworkServiceInfo(service)

		lines := strings.Split(info, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "IP address: ") {
				ipAddress := strings.TrimPrefix(line, "IP address: ")
				fmt.Printf("| %-15s | %-18s | %-10s | %-10s |\n", service, ipAddress, IP, City)
				fmt.Printf("--------------------------------------------------------------------\n")
			}
		}
	}

}

func CheckOS() {
	if runtime.GOOS != "darwin" {
		fmt.Println("[E] ONLY SUPPORTS MACOS, ABOUT TO EXIT......")
		os.Exit(1)
	}
}

func GetNetworkServices() []string {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	services := strings.Split(out.String(), "\n")
	var result []string
	for _, service := range services {
		if service != "" && !strings.HasPrefix(service, "*") {
			result = append(result, service)
		}
	}
	return result
}

func GetNetworkServiceInfo(service string) string {
	cmd := exec.Command("networksetup", "-getinfo", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	return out.String()
}

func GetPublic() (string, string) {
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", ""
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return "", ""
	}

	return ipInfo.IP, ipInfo.City
}
