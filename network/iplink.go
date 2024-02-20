package network

import (
	"net"
	"os"
	"strconv"

	"github.com/running910/gokit/logger"
	"github.com/running910/gokit/misc"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func DelNic(nic string) error {
	link, err := netlink.LinkByName(nic)
	if err != nil {
		logger.Errorf("netlink.LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	if err := netlink.LinkDel(link); err != nil {
		logger.Errorf("netlink.LinkDel() failed!, %s, %s", nic, err)
	}

	return nil
}

func AddVlanNic(parent string, nic string, vlanid uint32) error {
	link, err := netlink.LinkByName(parent)
	if err != nil {
		logger.Errorf("netlink.LinkByName() nic %s failed! reason: %s", parent, err)
		return err
	}

	newLink := &netlink.Vlan{
		LinkAttrs:    netlink.LinkAttrs{Name: nic, ParentIndex: link.Attrs().Index},
		VlanId:       int(vlanid),
		VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
	}

	if err := netlink.LinkAdd(newLink); err != nil {
		logger.Errorf("netlink.LinkAdd() parent:%s nic:%s failed! reason: %s", parent, nic, err)
		return err
	}

	return nil
}

func AddBridgeNic(nic string) error {

	link := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{Name: nic},
	}

	if err := netlink.LinkAdd(link); err != nil {
		logger.Errorf("netlink.LinkAdd() nic:%s failed! reason: %s", nic, err)
		return err
	}

	return nil
}

func AddMacvlanNic(parent string, nic string) error {
	link, err := netlink.LinkByName(parent)
	if err != nil {
		logger.Errorf("netlink.LinkByName() nic %s failed! reason: %s", parent, err)
		return err
	}

	newLink := &netlink.Macvlan{
		LinkAttrs: netlink.LinkAttrs{Name: nic, ParentIndex: link.Attrs().Index},
		Mode:      netlink.MACVLAN_MODE_BRIDGE,
	}

	if err := netlink.LinkAdd(newLink); err != nil {
		logger.Errorf("netlink.LinkAdd() parent:%s nic:%s failed! reason: %s", parent, nic, err)
		return err
	}

	return nil
}

func AddMacvlanNicBasedOnVlan(parent string, nic string, vlanid uint32) error {

	baseNic := parent

	// for example eth0 as parent and 100 as vlanid, need to make sure eth0.100 exist firstly
	if vlanid > 0 {
		baseNic = parent + "." + strconv.FormatUint(uint64(vlanid), 10)
		//	baseNic = vlanNic

		// create vlan interface if it does not exist
		if !CheckIfNicExist(baseNic) {

			// it must be wrong if logic nic already exist, need to delete it before create vlan nic
			if CheckIfNicExist(nic) {
				DelNic(nic)
			}

			logger.Info("vlan interface", baseNic, "not exist yet, now create it")
			AddVlanNic(parent, baseNic, vlanid)
			SetNsNicMacaddr(netns.None(), baseNic, misc.GenerateRandUnicastMacaddr())
			SetNicLinkUp(baseNic)
		}
	}

	// if logic nic does not exist, create it
	if !CheckIfNicExist(nic) {
		logger.Info("logic interface", nic, "not exist yet, now create it")
		AddMacvlanNic(baseNic, nic)
	} else {
		logger.Info("logic interface", nic, "exists already, do nothing")
	}

	SetNicLinkUp(nic)

	return nil
}

func SetNsNicMacaddr(ns netns.NsHandle, nic string, macaddr string) error {
	handle, err := netlink.NewHandleAt(ns)
	if err != nil {
		logger.Errorf("netlink.NewHandleAt() failed! reason:%s, ns:%d", err, ns)
		return err
	}
	defer handle.Delete()

	link, err := handle.LinkByName(nic)
	if err != nil {
		logger.Errorf("handle.LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	mac, err := net.ParseMAC(macaddr)
	if err != nil {
		logger.Errorf("net.ParseMAC() failed!, %s, %s", macaddr, err)
		return err
	}

	if err := handle.LinkSetHardwareAddr(link, []byte(mac)); err != nil {
		logger.Errorf("handle.LinkSetHardwareAddr() failed!, %s, %s", nic, mac, err)
		return err
	}

	return nil
}

func SetNicLinkUp(nic string) error {
	link, err := netlink.LinkByName(nic)
	if err != nil {
		logger.Errorf("netlink.LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		logger.Errorf("netlink.LinkSetUp() failed!, %s, %s", nic, err)
		return err
	}

	return nil
}

func SetNicLinkDown(nic string) error {
	link, err := netlink.LinkByName(nic)
	if err != nil {
		logger.Errorf("netlink.LinkByName() failed!, %s, %s", nic, err)
		return err
	}

	if err := netlink.LinkSetDown(link); err != nil {
		logger.Errorf("netlink.LinkSetDown() failed!, %s, %s", nic, err)
		return err
	}

	return nil
}

func CheckIfNicExist(nic string) bool {
	_, err := os.Stat("/sys/class/net/" + nic)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}
