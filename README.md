
- Requirements:
    
        go installed version: 1.14.3
        go get github.com/golang/glog
        go get gopkg.in/yaml.v2

- Build:

        go build
        
- Configuration:

        config.yml:
        
        networks:
          - network: "127.0.0.1/29"
            labels:
              networkzone: public
          - network: "10.0.0.5/27"
            labels:
              any: label
        concurrency: 11
        filesdpath: /file/path/file-sd.json
        port: 9090
        timeout: 5

    * to increase log level add parameter -stderrthreshold=INFO

- Run:

        $ ./host-discovery -c config.yml
        $ ./host-discovery -c config.yml -stderrthreshold=INFO


TODO: labels will be added in file_sd_configs.

