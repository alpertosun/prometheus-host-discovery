package main

import (
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func readYaml(filename string) (*Hosts,error) {
	var hosts Hosts
	yamlFile, err := os.Open(filename)
	if err != nil {
		glog.Error("cant open: ",filename," ",err)
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(yamlFile)
	err = yaml.Unmarshal(byteValue,&hosts)
	if err != nil {
		glog.Error("cant parse: ",filename," ",err)
		return nil,err
	}
	return &hosts,nil
}

func trimSpace(string2 string) string {
	trimmed := strings.TrimSpace(string2)
	return trimmed
}

func parseHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}