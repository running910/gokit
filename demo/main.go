package main

import (
	"github.com/running910/gokit/fs"
	"github.com/running910/gokit/logger"
	"github.com/running910/gokit/misc"
	"github.com/running910/gokit/network"
)

func main() {
	network.Hello()

	ipt, err := network.NewIptablesCtx()
	if err != nil {
		logger.Errorf("NewIptablesCtx() failed! reason:%s", err)
		return
	}

	ipt.EnsureChain(network.IpProtoV4, "filter", "hellochain")

	specs := []string{"-p", "udp", "-m", "multiport", "--dport", "100,200,300,147", "-j", "ACCEPT"}
	ipt.EnsureRuleInserted(network.IpProtoV4, "filter", "hellochain", specs...)
	if err != nil {
		logger.Errorf("EnsureRuleInserted() failed! reason:%s", err)
		return
	}

	specs = []string{"-j", "hellochain"}
	ipt.EnsureRuleInserted(network.IpProtoV4, "filter", "INPUT", specs...)
	if err != nil {
		logger.Errorf("EnsureRuleInserted() failed! reason:%s", err)
		return
	}

	specs = []string{"-j", "hellochain"}
	ipt.DeleteRule(network.IpProtoV4, "filter", "INPUT", specs...)

	ipt.DeleteChain(network.IpProtoV4, "filter", "hellochain")

	logger.Info(misc.CheckIfSliceContain([]int{1, 2, 4, 5}, 5))

	logger.Info(misc.CheckIfSliceContain([]int{1, 2, 4, 5}, 0))

	logger.Info(misc.CheckIfSliceContain([]string{"acc", "bbb", "eee"}, "a"))

	logger.Info(misc.CheckIfSliceContain([]string{"acc", "bbb", "eee"}, "bbb"))

	fs.Hello()
	misc.Hello()

	nic := "eth1"

	//misc.AttachPciDevDriver("0000:02:05.0", "e1000")
	//misc.AttachPciDevDriver("0000:06:00.1", "ixgbe")
	//os.Exit(-1)

	slot, err := misc.GetNicPciSlotId(nic)

	logger.Info(slot, err)

	driver, err := misc.GetNicDriver(nic)

	logger.Info(driver, err)

	vendor, err := misc.GetNicVendor(nic)

	logger.Info(vendor, err)

	device, err := misc.GetNicDevice(nic)

	logger.Info(device, err)

	err = misc.DetachPciDevDriver(slot, driver)
	logger.Info("err,", err)

	err = misc.AttachPciDevToVfioDriver(vendor, device)
	logger.Info("err,", err)

}
