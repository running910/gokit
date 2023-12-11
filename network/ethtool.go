package network

import (
	"github.com/running910/gokit/logger"
	"github.com/safchain/ethtool"
)

func EthtoolSetFeatureOnOff(nic string, feature string, value string) error {
	ethHandle, err := ethtool.NewEthtool()
	if err != nil {
		logger.Errorf("ethtool.NewEthtool() failed reason:%s", err)
		return err
	}
	defer ethHandle.Close()

	logger.Infof("ethtool -K %s %s %s", nic, feature, value)

	On := true
	if value == "off" {
		On = false
	}

	if err := ethHandle.Change(nic, map[string]bool{
		feature: On,
	}); err != nil {
		logger.Errorf("ethHandle.Change failed! reason:%s", err)
		return err
	}

	return nil
}
