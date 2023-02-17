package internal

import "net"

var localIP net.IP

func getLocalAddress() (net.IP, error) {
	con, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer con.Close()

	localAddress := con.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}
