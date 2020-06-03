package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func (c *Hosts) write(pingedHosts <-chan string)  {
	var fileSd []FileSD
	if _, err := os.Stat(c.FileSdPath); err == nil {
		jsonFile, err := os.Open(c.FileSdPath)
		if err != nil {
			log.Fatal("cant open",c.FileSdPath,err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		_=json.Unmarshal(byteValue,&fileSd)
	}
	for v := range pingedHosts {
		newElement := FileSD{
			Targets: []string{v},
			Labels: map[string]string{} ,
		}
		isIn := false
		for _,v := range fileSd {
			if reflect.DeepEqual(v.Targets,newElement.Targets) {
				isIn = true
				break
			}
		}
		if isIn {
			continue
		}
		fileSd = append(fileSd, newElement)
	}
	data, _ := json.Marshal(fileSd)
	_=ioutil.WriteFile(c.FileSdPath,data,0644)
	fmt.Println("Discovery finished. File writed at:",c.FileSdPath)
}

