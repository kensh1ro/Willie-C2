package scanutils

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type ScanResult struct {
	Port  string
	State string
}

var timeoutTCP time.Duration

func RunScan(host, ports, protocol string) string {

	timeoutTCP = time.Duration(500) * time.Millisecond

	portsList := getPortsList(ports)

	if portsList == nil {
		return "something went wrong"
	}

	var output = ""

	for p := range portsList {

		output += ScanPort(protocol, host, portsList[p])

	}

	// wait 1 sec more than timeout for finishing go routines
	time.Sleep(timeoutTCP + (1000 * time.Millisecond))
	return output
}

func getPortsList(port_var string) []string {
	// if port argument is like : 22,80,23
	if strings.Contains(port_var, ",") {
		ports_list := strings.Split(port_var, ",")

		for p := range ports_list {
			_, err_c := strconv.Atoi(ports_list[p])
			if err_c != nil {
				return nil
			}
		}

		return ports_list

	} else if strings.Contains(port_var, "-") {

		port_min_and_max := strings.Split(port_var, "-")

		port_min, err := strconv.Atoi(port_min_and_max[0])
		if err != nil {
			return nil

		}

		port_max, err := strconv.Atoi(port_min_and_max[1])
		if err != nil {
			return nil

		}

		var ports_temp_list []string

		for p_min := port_min; p_min <= port_max; p_min++ {
			port_str := strconv.Itoa(p_min)
			ports_temp_list = append(ports_temp_list, port_str)

		}

		return ports_temp_list

	}

	// if port is single number like : 80
	_, err := strconv.Atoi(port_var) // check if port is correct (int)
	if err != nil {
		return nil
	}

	return []string{port_var}

}

func ScanPort(protocol, hostname, port string) string {

	result := ScanResult{Port: port + string("/") + protocol + ", "}
	address := hostname + ":" + port
	conn, err := net.DialTimeout(protocol, address, 5*time.Second)

	if err != nil {
		result.State = "Closed\n"
	} else {
		defer conn.Close()
		result.State = "Open\n"
	}
	return result.Port + result.State
}
