
get_current_host_name() {
     echo $(hostname)
}

get_current_ip() {
     echo $(hostname -I | awk '{print $1}')
}

echo "Getting the host name and ip"
echo 
echo "Hostname:"
get_current_host_name

echo 
echo "Current Ip:"
get_current_ip
echo 