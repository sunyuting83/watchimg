package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host       string `yaml:"Host"`
	Code       string `yaml:"Code"`
	ComputName string `yaml:"ComputName"`
}

// CheckConfig check config
func CheckConfig(OS, CurrentPath string) (conf *Config, err error) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	if len(confYaml.Host) <= 0 {
		confYaml.Host = "http://192.168.0.178:13002/api/upimg"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.Code) <= 0 {
		mac, _ := GetMACAddressAndIPAddress()
		srcCode := md5.Sum([]byte(mac))
		code := fmt.Sprintf("%x", srcCode)
		confYaml.Code = string(code)
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.ComputName) <= 0 {
		mac, _ := GetMACAddressAndIPAddress()
		confYaml.ComputName = string(mac)
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	return confYaml, nil
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

// GetMACAddressAndIPAddress get MAC Address And IP Address
func GetMACAddressAndIPAddress() (mac string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}
	mac, macerr := "", errors.New("CantGetRealMACAddress")
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && (netInterfaces[i].Flags&net.FlagLoopback) == 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				ipnet, ok := address.(*net.IPNet)
				if ok && ipnet.IP.IsGlobalUnicast() {
					mac = netInterfaces[i].HardwareAddr.String()
					return mac, nil
				}
			}
		}
	}
	return mac, macerr
}
