package main

import (
	"github.com/golang/glog"
	"net"
	"sync"
	"time"
)


type Hosts struct {
	Networks []struct {
		Network string            `yaml:"network"`
		Labels  map[string]string `yaml:"labels"`
	} `yaml:"networks"`
	Concurrency int           `yaml:"concurrency"`
	Port        string           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	FileSdPath  string        `yaml:"filesdpath"`
}

type FileSD struct {
	Targets []string `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func (c *Hosts) receiveHosts() ([]string,error) {
	var hostList []string

	for _,val := range c.Networks {
		hosts, err := parseHosts(val.Network)
		if err != nil {
			glog.Error(err)
			return nil,err
		}
		for _,host := range hosts {
			hostList = append(hostList, host)
		}
	}
	glog.Info("Total number of hosts to discover: ", len(hostList))
	return hostList,nil
}

func (c *Hosts) lookup(hostList []string) <-chan string {
	var hostChannel = make(chan string)
	var wg = sync.WaitGroup{}

	//blocking pings
	hostNumber := len(hostList)
	numberOfJobs := hostNumber / c.Concurrency
	glog.Info("Worker pool(job*work): ",numberOfJobs," * ",c.Concurrency)

	job := 0
	go func() {
		for j:=0;j<numberOfJobs;j++  {
			for i:=0;i<c.Concurrency;i++ {
				glog.Info(job," - Starting dial-up: ",hostList[job]+":"+c.Port)
				wg.Add(1)
				IsOpen(hostList[job],c.Port,c.Timeout,hostChannel,&wg)
				job ++
			}
			wg.Wait()
		}
		if hostNumber % c.Concurrency != 0 {
			for j := job;j<hostNumber;j++ {
				wg.Add(1)
				glog.Info(job," - Starting dial-up: ",hostList[job]+":"+c.Port)
				IsOpen(hostList[job],c.Port,c.Timeout,hostChannel,&wg)
				job ++
			}
			wg.Wait()
		}
		close(hostChannel)
	}()


	return hostChannel
}

func IsOpen(ip string,port string,timeout time.Duration,hostChannel chan string,group *sync.WaitGroup)  {
	if ip == "" {return}
	go func() {
		defer group.Done()
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip,port), timeout*time.Second)
		if conn != nil {
			defer conn.Close()
			glog.Info("Connected to: ",net.JoinHostPort(ip,port))
			hostChannel <- net.JoinHostPort(ip,port)
		}
		if err != nil {
			glog.Warning("Can not to: ",net.JoinHostPort(ip,port),"- error:",err)
			return
		}
	}()
}







