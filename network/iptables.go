package network

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/running910/gokit/logger"
)

type IptablesCtx struct {
	ip4t *iptables.IPTables
	ip6t *iptables.IPTables
}

func NewIptablesCtx() (*IptablesCtx, error) {

	ip4t, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		logger.Errorf("NewWithProtocol() failed with proto ipv4! reason:%s", err)
		return nil, err
	}

	ip6t, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		logger.Errorf("NewWithProtocol() failed with proto ipv4! reason:%s", err)
		return nil, err
	}

	return &IptablesCtx{ip4t: ip4t, ip6t: ip6t}, nil

}

func (i *IptablesCtx) Release() {

}

func (i *IptablesCtx) EnsureChain(proto IpProto, table, chain string) error {

	ipt := i.ip4t
	if proto == IpProtoV6 {
		ipt = i.ip6t
	}

	exist, err := ipt.ChainExists(table, chain)
	if err != nil {
		logger.Errorf("ChainExists for existing chain failed: %v\n", err)
		return err
	} else if !exist {
		logger.Errorf("ChainExists doesn't find existing chain")

		err = ipt.ClearChain(table, chain)
		if err != nil {
			logger.Errorf("ClearChain (of empty) failed: %v\n", err)
			return err
		}
	}

	return nil

}

func (i *IptablesCtx) EnsureRuleAppended(proto IpProto, table, chain string, specs ...string) error {

	ipt := i.ip4t
	if proto == IpProtoV6 {
		ipt = i.ip6t
	}

	exist, err := ipt.Exists(table, chain, specs...)
	if err != nil {
		logger.Errorf("Exists for existing chain failed: %v", err)
		return err
	} else if !exist {
		logger.Errorf("Exists doesn't find existing rule")

		err = ipt.Append(table, chain, specs...)
		//err = ipt.Insert(table, chain, 1, specs...)
		if err != nil {
			logger.Errorf("Append failed: %v\n", err)
			return err
		}
	}

	return nil
}

func (i *IptablesCtx) EnsureRuleInserted(proto IpProto, table, chain string, specs ...string) error {

	ipt := i.ip4t
	if proto == IpProtoV6 {
		ipt = i.ip6t
	}

	exist, err := ipt.Exists(table, chain, specs...)
	if err != nil {
		logger.Errorf("Exists for existing chain failed: %v", err)
		return err
	} else if !exist {
		logger.Errorf("Exists doesn't find existing rule")

		//err = ipt.Append(table, chain, specs...)
		err = ipt.Insert(table, chain, 1, specs...)
		if err != nil {
			logger.Errorf("Append failed: %v\n", err)
			return err
		}
	}

	return nil
}
