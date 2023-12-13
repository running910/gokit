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

	fs.Hello()
	misc.Hello()
}
