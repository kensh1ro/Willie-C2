package netutils

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func GetUrlBytes(url string) (data []byte, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func DownloadURL(filepath string, url string) error {

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	return err
}

func PrivateIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err.Error()
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func Ipconfig() (string, error) {
	var out strings.Builder
	b := make([]byte, 1)
	l := uint32(len(b))
	aList := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err := syscall.GetAdaptersInfo(aList, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		aList = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(aList, &l)
	}
	if err != nil {
		return "", os.NewSyscallError("GetAdaptersInfo", err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, ifi := range ifaces {
		for pAdapt := aList; pAdapt != nil; pAdapt = pAdapt.Next {
			index := pAdapt.Index
			out.WriteByte(10)
			out.WriteString("Adpater:\t")
			out.WriteString(string(pAdapt.Description[:]))
			out.WriteByte(10)
			out.WriteString("\tPhysical Address:\t")
			out.WriteString(ifi.HardwareAddr.String())
			out.WriteByte(10)
			if ifi.Index == int(index) {
				pAddrStr := &pAdapt.IpAddressList
				for ; pAddrStr != nil; pAddrStr = pAddrStr.Next {
					out.WriteString("\tIP Address:\t")
					out.WriteString(string(pAddrStr.IpAddress.String[:]))
					out.WriteByte(10)
					out.WriteString("\tSubnet Mask:\t ")
					out.WriteString(string(pAddrStr.IpMask.String[:]))
					out.WriteByte(10)
					out.WriteString("\tDefault Gateway:\t")
					out.WriteString(string(pAdapt.GatewayList.IpAddress.String[:]))
					out.WriteByte(10)
				}

				pAddrStr = &pAdapt.GatewayList
				for ; pAddrStr != nil; pAddrStr = pAddrStr.Next {
					out.WriteString("\tGateway:\t")
					out.WriteString(string(pAddrStr.IpAddress.String[:]))
					out.WriteByte(10)
				}
				if pAdapt.DhcpEnabled != 0 {
					dhcpServers := &pAdapt.DhcpServer
					for ; dhcpServers != nil; dhcpServers = dhcpServers.Next {
						out.WriteString("\tDHCP Server:\t")
						out.WriteString(string(dhcpServers.IpAddress.String[:]))
						out.WriteByte(10)
					}
				} else {
					out.WriteString("\tDHCP Disabled\n")
				}
			}
		}
	}
	return out.String(), err
}
