
### Build

Linux: 

`env GOOS=linux GOARCH=amd64 go build -o link-monitor-linux`



## Todo

- implement logging instead of fmt.Print
- add broker config using the conf file
- publish to correct routing key with kit Id replacemnt
- persist counters in cache dir
- all the devops
- expand message struct to support a full stat output per net interface
- proper cleanup - close channels first