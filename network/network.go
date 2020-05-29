package network

import (
	"net"
)

func CIDRToIPRange(cidr string) (ips []string, err error){
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	ip := ipNet.IP.Mask(ipNet.Mask)
	for {
		if !ipNet.Contains(ip) {
			break
		}
		ips = append(ips, ip.String())
		for i := len(ip)-1; i >= 0; i-- {
			ip[i]++
			if ip[i] > 0 {
				break
			}
		}
	}
	if len(ips) == 0 {
		return
	}
	ips = ips[1:len(ips)-1]
	return
}
