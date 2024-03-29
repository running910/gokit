package network

import (
	"testing"
)

func TestGetNicNetlinkIndex(t *testing.T) {

	index := GetNicNetlinkIndex("ens33")

	if index <= 0 {
		t.Fatalf("get ens33 index:%d less than or equal 0", index)
	}

}
