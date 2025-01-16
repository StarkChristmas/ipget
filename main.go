package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/liushuochen/gotable"
	"io/ioutil"
	"log"
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

type NetworkSetup struct {
	NetworkService  string `json:"network_service"`
	ConfigureMethod string `json:"configure_method"`
	IPAddress       string `json:"ip_address"`
	SubnetMask      string `json:"subnet_mask"`
	Router          string `json:"router"`
	IPv6Method      string `json:"ipv6_method"`
	IPv6Address     string `json:"ipv6_address"`
	IPv6Route       string `json:"ipv6_route"`
	EthernetAddress string `json:"ethernet_address"`
}

type NetworkDetailsInfo struct {
	PublicIP     string         `json:"public_ip"`
	City         string         `json:"city"`
	NetworkSetup []NetworkSetup `json:"network_service"`
}

func main() {
	CheckOS()
	IP, City := GetPublic()
	ActiveNetworkInterface := GetNetworkServices()
	var strSlice []string
	for _, activeNetworkInterface := range ActiveNetworkInterface {
		info := GetNetworkServiceInfo(activeNetworkInterface)
		NetworkInfo := "NetworkName: " + activeNetworkInterface + "\n" + info
		strSlice = append(strSlice, NetworkInfo)

	}
	result := strings.Join(strSlice, " ")

	setups := parseNetworkData(result)

	jsonData, err := json.MarshalIndent(setups, "", "  ")
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
		return
	}

	var networkInterfaces []NetworkSetup

	if err := json.Unmarshal(jsonData, &networkInterfaces); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	updatedConfig := NetworkDetailsInfo{
		PublicIP:     IP,
		City:         City,
		NetworkSetup: setups,
	}

	table, err := gotable.Create("Network Service", "Local IPv4 Address", "Public IP", "City")
	if err != nil {
		log.Fatalf("Create table failed: %v", err)
	}
	for _, ni := range updatedConfig.NetworkSetup {
		if ni.EthernetAddress != "(null)" && ni.EthernetAddress != "" {
			table.AddRow([]string{
				ni.NetworkService,
				ni.IPAddress,
				IP,
				City,
			})
		}

	}
	fmt.Println(table)
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
	ActiveNetworkInterface := strings.Split(out.String(), "\n")
	var result []string
	for _, service := range ActiveNetworkInterface {
		if service != "" && !strings.HasPrefix(service, "*") {
			result = append(result, service)
		}
	}
	currentResult := filterString(result, "An asterisk (*) denotes that a network service is disabled.")
	return currentResult
}

func GetNetworkServiceInfo(service string) string {
	cmd := exec.Command("networksetup", "-getinfo", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	output := out.String()

	excludeUseless := strings.Replace(output, "An asterisk (*) denotes that a network service is disabled. is not a recognized network service.", "", -1)
	excludeUseless1 := strings.Replace(excludeUseless, "** Error: The parameters were not valid.", "", -1)
	return excludeUseless1
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

func parseNetworkData(data string) []NetworkSetup {
	var networks []NetworkSetup
	lines := strings.Split(data, "\n")
	var currentNetwork NetworkSetup

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "NetworkName:") {
			if currentNetwork.NetworkService != "" {
				networks = append(networks, currentNetwork)
			}
			currentNetwork = NetworkSetup{}
			currentNetwork.NetworkService = strings.TrimSpace(line[len("NetworkName: "):])
		} else if line == "DHCP Configuration" || line == "Manual Configuration" {
			currentNetwork.ConfigureMethod = line
		} else if strings.HasPrefix(line, "IP address:") {
			currentNetwork.IPAddress = strings.TrimSpace(line[len("IP address: "):])
		} else if strings.HasPrefix(line, "Subnet mask:") {
			currentNetwork.SubnetMask = strings.TrimSpace(line[len("Subnet mask: "):])
		} else if strings.HasPrefix(line, "Router:") {
			currentNetwork.Router = strings.TrimSpace(line[len("Router: "):])
		} else if strings.HasPrefix(line, "IPv6:") {
			currentNetwork.IPv6Method = strings.TrimSpace(line[len("IPv6: "):])
		} else if strings.HasPrefix(line, "IPv6 IP address:") {
			currentNetwork.IPv6Address = strings.TrimSpace(line[len("IPv6 IP address: "):])
		} else if strings.HasPrefix(line, "IPv6 Router:") {
			currentNetwork.IPv6Route = strings.TrimSpace(line[len("IPv6 Router: "):])
		} else if strings.HasPrefix(line, "Ethernet Address:") {
			currentNetwork.EthernetAddress = strings.TrimSpace(line[len("Ethernet Address: "):])
		} else if strings.HasPrefix(line, "Wi-Fi ID:") {
			currentNetwork.EthernetAddress = strings.TrimSpace(line[len("Wi-Fi ID: "):])
		}
	}
	if currentNetwork.NetworkService != "" {
		networks = append(networks, currentNetwork)
	}

	return networks
}

func filterString(slice []string, toRemove string) []string {
	var filtered []string
	for _, str := range slice {
		if str != toRemove {
			filtered = append(filtered, str)
		}
	}
	return filtered
}
