# bgptracer
BGP path tracer for AWS direct connect

# Prerequisites

- You must have go installed for building the program
- You must have nmap installed on the machine where you will run this, as this program uses nmap


# Install / Build

- clone the repo
- cd bgptracer
- Run go get ./... to install the dependencies



# Usage

- Modify the code to setup your IP address map
- Setup your slack webhook and update the URL in the code. (This is optional. If slack post is not needed this can be removed)

- go build bgptracer.go
- Run bgptracer in your on-prem network and use a target host/port in the AWS cloud which is reachable :

$ ./bgptracer <targetport> <targethost> <detailed_log_file_path> <summary_log_file_path> <delay_in_seconds>
