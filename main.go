package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/liushuochen/gotable"
)

// IPInfo 存储从ipinfo.io获取的公共IP信息
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

// NetworkSetup 存储网络接口的详细信息
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

// 错误信息结构体
type ErrorMessages struct {
	OSNotSupported       string
	RequestFailed        string
	ReadResponseFailed   string
	ParseJSONFailed      string
	CreateTableFailed    string
	ExecuteCommandFailed string
	GetNetworkInfoFailed string
}

// 根据系统语言环境获取错误信息
func getErrorMessages() ErrorMessages {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "zh") {
		return ErrorMessages{
			OSNotSupported:       "仅支持macOS系统",
			RequestFailed:        "请求失败: %w",
			ReadResponseFailed:   "读取响应失败: %w",
			ParseJSONFailed:      "解析JSON失败: %w",
			CreateTableFailed:    "创建表格失败: %w",
			ExecuteCommandFailed: "执行命令失败: %w",
			GetNetworkInfoFailed: "获取网络信息失败: %w",
		}
	}
	return ErrorMessages{
		OSNotSupported:       "Only macOS is supported",
		RequestFailed:        "Request failed: %w",
		ReadResponseFailed:   "Failed to read response: %w",
		ParseJSONFailed:      "Failed to parse JSON: %w",
		CreateTableFailed:    "Failed to create table: %w",
		ExecuteCommandFailed: "Failed to execute command: %w",
		GetNetworkInfoFailed: "Failed to get network info: %w",
	}
}

func main() {
	if err := run(); err != nil {
		lang := os.Getenv("LANG")
		if strings.HasPrefix(lang, "zh") {
			log.Fatalf("程序执行失败: %v", err)
		} else {
			log.Fatalf("Program execution failed: %v", err)
		}
	}
}

// run 是主程序的入口函数，处理所有主要逻辑
func run() error {
	errors := getErrorMessages()

	if err := checkOS(); err != nil {
		return fmt.Errorf(errors.OSNotSupported)
	}

	ip, city, err := getPublicIP()
	if err != nil {
		return fmt.Errorf(errors.RequestFailed, err)
	}

	networkInterfaces, err := getActiveNetworkInterfaces()
	if err != nil {
		return fmt.Errorf(errors.GetNetworkInfoFailed, err)
	}

	if err := displayNetworkInfo(networkInterfaces, ip, city); err != nil {
		return fmt.Errorf(errors.CreateTableFailed, err)
	}

	return nil
}

// checkOS 检查操作系统是否为macOS
func checkOS() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf(getErrorMessages().OSNotSupported)
	}
	return nil
}

// getPublicIP 获取公共IP地址和城市信息
func getPublicIP() (string, string, error) {
	errors := getErrorMessages()
	resp, err := http.Get("http://ipinfo.io")
	if err != nil {
		return "", "", fmt.Errorf(errors.RequestFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf(errors.ReadResponseFailed, err)
	}

	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", "", fmt.Errorf(errors.ParseJSONFailed, err)
	}

	return ipInfo.IP, ipInfo.City, nil
}

// getActiveNetworkInterfaces 获取活动的网络接口信息
func getActiveNetworkInterfaces() ([]NetworkSetup, error) {
	services, err := getNetworkServices()
	if err != nil {
		return nil, err
	}

	var networkInfo []string
	for _, service := range services {
		info, err := getNetworkServiceInfo(service)
		if err != nil {
			return nil, err
		}
		networkInfo = append(networkInfo, fmt.Sprintf("NetworkName: %s\n%s", service, info))
	}

	return parseNetworkData(strings.Join(networkInfo, " ")), nil
}

// getNetworkServices 获取所有网络服务
func getNetworkServices() ([]string, error) {
	errors := getErrorMessages()
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(errors.ExecuteCommandFailed, err)
	}

	services := strings.Split(out.String(), "\n")
	var result []string
	for _, service := range services {
		if service != "" && !strings.HasPrefix(service, "*") {
			result = append(result, service)
		}
	}

	return filterString(result, "An asterisk (*) denotes that a network service is disabled."), nil
}

// getNetworkServiceInfo 获取指定网络服务的详细信息
func getNetworkServiceInfo(service string) (string, error) {
	errors := getErrorMessages()
	cmd := exec.Command("networksetup", "-getinfo", service)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(errors.ExecuteCommandFailed, err)
	}

	output := out.String()
	output = strings.ReplaceAll(output, "An asterisk (*) denotes that a network service is disabled. is not a recognized network service.", "")
	output = strings.ReplaceAll(output, "** Error: The parameters were not valid.", "")

	return output, nil
}

// displayNetworkInfo 显示网络信息表格
func displayNetworkInfo(networks []NetworkSetup, publicIP, city string) error {
	errors := getErrorMessages()
	table, err := gotable.Create("Network Service", "Local IPv4 Address", "Public IP", "City")
	if err != nil {
		return fmt.Errorf(errors.CreateTableFailed, err)
	}

	for _, ni := range networks {
		if ni.EthernetAddress != "(null)" && ni.EthernetAddress != "" && ni.IPAddress != "" {
			table.AddRow([]string{
				ni.NetworkService,
				ni.IPAddress,
				publicIP,
				city,
			})
		}
	}

	fmt.Println(table)
	return nil
}

// parseNetworkData 解析网络数据字符串为NetworkSetup结构体切片
func parseNetworkData(data string) []NetworkSetup {
	var networks []NetworkSetup
	lines := strings.Split(data, "\n")
	var currentNetwork NetworkSetup

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "NetworkName:"):
			if currentNetwork.NetworkService != "" {
				networks = append(networks, currentNetwork)
			}
			currentNetwork = NetworkSetup{}
			currentNetwork.NetworkService = strings.TrimSpace(line[len("NetworkName: "):])
		case line == "DHCP Configuration" || line == "Manual Configuration":
			currentNetwork.ConfigureMethod = line
		case strings.HasPrefix(line, "IP address:"):
			currentNetwork.IPAddress = strings.TrimSpace(line[len("IP address: "):])
		case strings.HasPrefix(line, "Subnet mask:"):
			currentNetwork.SubnetMask = strings.TrimSpace(line[len("Subnet mask: "):])
		case strings.HasPrefix(line, "Router:"):
			currentNetwork.Router = strings.TrimSpace(line[len("Router: "):])
		case strings.HasPrefix(line, "IPv6:"):
			currentNetwork.IPv6Method = strings.TrimSpace(line[len("IPv6: "):])
		case strings.HasPrefix(line, "IPv6 IP address:"):
			currentNetwork.IPv6Address = strings.TrimSpace(line[len("IPv6 IP address: "):])
		case strings.HasPrefix(line, "IPv6 Router:"):
			currentNetwork.IPv6Route = strings.TrimSpace(line[len("IPv6 Router: "):])
		case strings.HasPrefix(line, "Ethernet Address:"):
			currentNetwork.EthernetAddress = strings.TrimSpace(line[len("Ethernet Address: "):])
		case strings.HasPrefix(line, "Wi-Fi ID:"):
			currentNetwork.EthernetAddress = strings.TrimSpace(line[len("Wi-Fi ID: "):])
		}
	}

	if currentNetwork.NetworkService != "" {
		networks = append(networks, currentNetwork)
	}

	return networks
}

// filterString 从字符串切片中过滤掉指定的字符串
func filterString(slice []string, toRemove string) []string {
	var filtered []string
	for _, str := range slice {
		if str != toRemove {
			filtered = append(filtered, str)
		}
	}
	return filtered
}
