package network

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/running910/gokit/logger"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func ParseCIDR(ipaddr string, netmask string) (net.IP, *net.IPNet, error) {
	mask := net.IPMask(net.ParseIP(netmask).To4())
	merge := &net.IPNet{IP: net.ParseIP(ipaddr), Mask: mask}

	ip, network, err := net.ParseCIDR(merge.String())
	if err != nil {
		logger.Errorf("net.ParseCIDR() failed!, ipaddr:%s mask:%s reaon:%s", ipaddr, netmask, err)

		return net.IP{}, &net.IPNet{}, err
	}

	return ip, network, nil
}

func GetNsNicDefaultGateway(ns netns.NsHandle, nic string) (string, error) {
	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return "", err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return "", err
	}

	routes, err := handle.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		logger.Errorf("netlink.RouteList() failed!, reason: %s", err)
		return "", err
	}

	for _, route := range routes {
		if route.Table == 254 && route.Gw != nil && route.LinkIndex == link.Attrs().Index {
			return route.Gw.String(), nil
		}
	}
	return "", errors.New("get ns default gateway failed!")

}

func GetNsDefaultGateway(ns netns.NsHandle) (string, error) {
	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return "", err
	}
	defer handle.Delete()

	routes, err := handle.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		logger.Errorf("netlink.RouteList() failed!, reason: %s", err)
		return "", err
	}

	for _, route := range routes {
		if route.Table == 254 && route.Gw != nil {
			return route.Gw.String(), nil
		}
	}
	return "", errors.New("get ns default gateway failed!")

}

func CleanNsNicIpaddrv6Info(ns netns.NsHandle, nic string) error {

	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	ips, err := handle.AddrList(link, syscall.AF_INET6)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", nic, err)
		return err
	}

	for _, ip := range ips {
		if netlink.Scope(ip.Scope) == netlink.SCOPE_LINK {
			//logger.Info("ignore link local address:", ip.IP.String())
			continue
		}

		logger.Errorf("now remove nic %s ip address %s", nic, ip.String())

		handle.AddrDel(link, &ip)
	}

	return nil
}

func CleanNsNicIpaddrInfo(ns netns.NsHandle, nic string) error {

	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	ips, err := handle.AddrList(link, syscall.AF_INET)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", nic, err)
		return err
	}

	for _, ip := range ips {
		logger.Errorf("now remove nic %s ip address %s", nic, ip.String())
		handle.AddrDel(link, &ip)
	}

	return nil
}

func DelNsNicIpaddrInfo(ns netns.NsHandle, nic string, ipaddr string) error {

	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	ips, err := handle.AddrList(link, syscall.AF_UNSPEC)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", nic, err)
		return err
	}

	for _, ip := range ips {

		if ip.IP.String() == ipaddr {
			logger.Errorf("now remove nic %s ip address %s", nic, ip.IP.String())
			handle.AddrDel(link, &ip)
		}
	}

	return nil
}

func GetNsNicFirstIpaddrAndNetmask(ns netns.NsHandle, nic string) (string, string, string, error) {

	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return "", "", "", err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		//logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return "", "", "", err
	}

	ips, err := handle.AddrList(link, syscall.AF_INET)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", nic, err)
		return "", "", "", err
	}

	if len(ips) == 0 {
		return "", "", "", errors.New(fmt.Sprintf("nic %s no ipv4 address found", nic))
	}

	// if nic is ppp interface, there would be an peer address
	var peer string
	if ips[0].Peer != nil {
		peer = ips[0].Peer.IP.String()
	}

	return ips[0].IP.String(), net.IP(ips[0].Mask).String(), peer, nil
}

func GetNsNicAllIpaddrInfo(ns netns.NsHandle, nic string) ([]string, error) {

	ipaddrs := make([]string, 0)

	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return ipaddrs, err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		//logger.Errorf("LinkByName() failed!, %s, %s", nic, err)
		return ipaddrs, err
	}

	ips, err := handle.AddrList(link, syscall.AF_UNSPEC)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", nic, err)
		return ipaddrs, err
	}

	if len(ips) == 0 {
		return ipaddrs, errors.New(fmt.Sprintf("nic %s no ip address found", nic))
	}

	for _, ip := range ips {
		ipaddrs = append(ipaddrs, ip.IP.String())
	}

	return ipaddrs, nil
}

func GetNicFirstIpaddrAndNetmask(nic string) (string, string, string, error) {
	link, err := netlink.LinkByName(nic)
	if err != nil {
		logger.Errorf("netlink.LinkByName() failed!, %s, %s", nic, err)
		return "", "", "", err
	}

	ips, err := netlink.AddrList(link, syscall.AF_INET)
	if err != nil {
		logger.Errorf("netlink.AddrList() failed, %s, %s", "ens33", err)
		return "", "", "", err
	}

	if len(ips) == 0 {
		return "", "", "", errors.New(fmt.Sprintf("nic %s no ipv4 address found", nic))
	}

	// if nic is ppp interface, there would be an peer address
	var peer string
	if ips[0].Peer != nil {
		peer = ips[0].Peer.IP.String()
	}

	return ips[0].IP.String(), net.IP(ips[0].Mask).String(), peer, nil
}
