package network

import (
	"fmt"

	"github.com/running910/gokit/logger"
	"github.com/vishvananda/netlink"
)

func Hello() {
	fmt.Println("this is from gokit network package.")
}

func GetNicNetlinkIndex(nic string) int {
	link, err := netlink.LinkByName(nic)
	if err != nil {
		logger.Errorf("netlink.LinkByName() for nic:%s failed! reason:%s", nic, err)
		return 0
	}
	return link.Attrs().Index
}
