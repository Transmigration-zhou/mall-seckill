package common

import (
	"errors"
	"net"
)

// GetEntranceIp 获取本机ip
func GetEntranceIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		////检查Ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("获取地址异常")
}
