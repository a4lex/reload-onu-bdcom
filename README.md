# reload_onu_bdcom
Allow reload remotely ONU via Telnet on BDCom. 

##### Compile  

>https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04. 

env GOOS=linux GOARCH=amd64 go build  reload_onu.go

##### Usage of reload_onu:
- reload_onu -h
- reload_onu -host 10.1.1.1 -password MyPass -enable myEnablePass -onu-mac e0e8.e611.ed5e
- reload_onu -host 10.2.2.2 -password MyPass -enable myEnablePass -onu-path ./myOnuList
- reload_onu -host 10.3.3.3 -password MyPass  -onu-path ./myOnuList -log-level 63
