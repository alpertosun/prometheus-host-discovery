
- Requirements:
    
        go get 

- Build:

        go build
        
- Configuration:

        networks:
          # Network addresses for scanning. example: 127.0.0.1/24
          - network: <ip/net>
        # Number of simultaneous connection
        concurrency: <duration>
        port: 
          - <port>
        # Connection timeout for every attempt
        timeout: <duration>

- Run:

        $ ./host-discovery -c config.yml -f results.json
        
        
* to increase log level add parameter -stderrthreshold=INFO

        $ ./host-discovery -c config.yml -stderrthreshold=INFO
        
        
        
* give results.json to prometheus as file_sd_configs


TODO: labels will be added in file_sd_configs.

