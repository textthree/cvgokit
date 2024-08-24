package netkit

import (
	"encoding/binary"
	"net"
	"os"
	"strings"
)

// 获取本机名称
func Gethostname() (string, error) {
	return os.Hostname()
}

// 用域名或主机名获取IP地址，用于本地主机的标准主机名
func Gethostbyname(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", nil
	}
	return "", err
}

// 获取互联网主机名对应的 IPv4 地址列表，即获取同ip网站
func Gethostbynamel(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		var ipstrs []string
		for _, v := range ips {
			if v.To4() != nil {
				ipstrs = append(ipstrs, v.String())
			}
		}
		return ipstrs, nil
	}
	return nil, err
}

// 通过一个IPv4的地址来获取主机名
func Gethostbyaddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}

// 将IPV4 的字符串互联网协议转换成长整型数字
func Ip2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// 将长整型转化为字符串形式带点的互联网标准格式地址(IPV4)
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}
