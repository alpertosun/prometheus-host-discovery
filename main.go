package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string] -c config.yml\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	configYml := flag.String("c","config.yml","config file location")
	flag.Usage = usage
	flag.Parse()

	configData,err := readYaml(*configYml)
	if err !=nil {
		return
	}

	hostList, err := configData.receiveHosts()
	if err !=nil {
		return
	}

	hostsChannel := configData.lookup(hostList)
	configData.write(hostsChannel)
}
