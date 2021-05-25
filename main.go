package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/golang/glog"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string] -c config.yml\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	configLocation := flag.String("c", "config.yml", "location of config file")
	resultLocation := flag.String("f", "results.json", "location of config file")
	flag.Usage = usage
	flag.Parse()

	sdConfig, _ := readYaml(*configLocation)

	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	var hostList []string

	for _, networks := range sdConfig.Networks {
		hosts, _ := receiveHosts(networks.Network)
		hostList = append(hostList, hosts...)
	}

	count := len(hostList)
	bar := pb.StartNew(count)
	bar.SetWriter(os.Stdout)

	hostChan := make(chan string)
	for _, i := range hostList {
		semaphore <- struct{}{}
		wg.Add(1)
		for _, port := range sdConfig.Port {
			go IsOpen(i, strconv.Itoa(port), time.Duration(sdConfig.Concurrency), hostChan, semaphore, &wg, bar)
		}
	}
	go func() {
		wg.Wait()
		close(semaphore)
		close(hostChan)
	}()
	result := ParseSDConfig(hostChan)
	bar.Finish()
	fmt.Println(result)
	_ = ioutil.WriteFile(*resultLocation, []byte(result), 0644)
}

func ParseSDConfig(hosts chan string) string {
	g := struct {
		Targets []string `json:"targets"`
		Labels  string   `json:"-"`
	}{}
	for i := range hosts {
		g.Targets = append(g.Targets, i)
	}
	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
func IsOpen(ip string, port string, timeout time.Duration, hostChannel chan string, semaphore chan struct{}, wg *sync.WaitGroup, pb *pb.ProgressBar) {
	defer func() {
		<-semaphore
		wg.Done()
		pb.Add(1)
	}()
	if ip == "" {
		return
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout*time.Second)
	if conn != nil {
		_ = conn.Close()
		glog.Info("Connected to: ", net.JoinHostPort(ip, port))
		hostChannel <- net.JoinHostPort(ip, port)
	}
	if err != nil {
		glog.Warning("Can not connect to: ", net.JoinHostPort(ip, port), "- error:", err)
		return
	}
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
func receiveHosts(ipNet string) ([]string, error) {
	var hostList []string

	hosts, err := parseHosts(ipNet)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	for _, host := range hosts {
		hostList = append(hostList, host)
	}

	glog.Info("Total number of hosts to discover: ", len(hostList))
	return hostList, nil
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
