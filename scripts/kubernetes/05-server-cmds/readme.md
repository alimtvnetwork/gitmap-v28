# Command

```bash
./02-run-cmd-v2.sh all "hostnamectl hostname && echo $(hostname -I | awk '{print $1}')" 
./02-run-cmd-v2.sh all "ls -la /tmp"
./05-run-script-nodes.sh all ./06-echo.sh
```
## Get all the Ip and Host name run as below

```bash
./02-run-cmd-v2.sh all "sudo github/kubernetes-training-v1/ping.sh"
```

## Getting ip

```bash
# Get the IP address of the machine
ip_address=$(hostname -I | awk '{print $1}')

# Print the IP address and a greeting
echo "Hello from Ubuntu! Your IP address is $ip_address"
```