package misc

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	lg "github.com/running910/gokit/logger"
)

func Hello() {
	fmt.Println("this is from gokit misc package.")
}

// status orignal driverdetached vfioattached unknown
type pci_nic struct {
	Name   string
	Driver string
	SlotId string
	Vendor string
	Device string
	Status int
}

func (p *pci_nic) UpdateStatus() {

}

func GetNicVendor(nic string) (string, error) {
	file := fmt.Sprintf("/sys/class/net/%s/device/vendor", nic)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func GetNicDevice(nic string) (string, error) {
	file := fmt.Sprintf("/sys/class/net/%s/device/device", nic)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func GetNicDriver(nic string) (string, error) {

	// realpath /sys/class/net/ens37/device/driver/module/
	moduleFile := fmt.Sprintf("/sys/class/net/%s/device/driver/module", nic)

	devicePath, err := os.Readlink(moduleFile)
	if err != nil {
		return "", err
	}

	lg.Info("moduleFile: ", moduleFile, "drivername", devicePath)

	tokens := strings.Split(devicePath, "/")
	if len(tokens) < 1 {
		return "", fmt.Errorf("nic device path format error")
	}

	return tokens[len(tokens)-1], nil

}

func GetNicPciSlotId(nic string) (string, error) {

	interfacePath := filepath.Join("/sys/class/net", nic)
	devicePath, err := os.Readlink(interfacePath)
	if err != nil {
		return "", err
	}

	tokens := strings.Split(devicePath, "/")
	if len(tokens) < 3 {
		return "", fmt.Errorf("nic device path format error")
	}
	pciSlot := tokens[len(tokens)-3]

	return pciSlot, nil
}

func GetPciDevDriver(pciSlotId string) (string, error) {
	return "", nil
}

func DetachPciDevDriver(pciSlotId string, driver string) error {
	unbindFile := fmt.Sprintf("/sys/bus/pci/drivers/%s/unbind", driver)

	return ioutil.WriteFile(unbindFile, []byte(pciSlotId), os.ModeExclusive)
}

func AttachPciDevDriver(pciSlotId string, driver string) error {
	bindFile := fmt.Sprintf("/sys/bus/pci/drivers/%s/bind", driver)

	return ioutil.WriteFile(bindFile, []byte(pciSlotId), os.ModeExclusive)
}

func AttachPciDevToVfioDriver(vendor string, device string) error {
	bindFile := fmt.Sprintf("/sys/bus/pci/drivers/vfio-pci/new_id")

	content := fmt.Sprintf("%s %s", vendor, device)

	return ioutil.WriteFile(bindFile, []byte(content), os.ModeExclusive)
}

func DetachPciDevToVfioDriver(pciSlotId string) error {
	bindFile := fmt.Sprintf("/sys/bus/pci/drivers/vfio-pci/unbind")

	content := fmt.Sprintf("%s", pciSlotId)

	return ioutil.WriteFile(bindFile, []byte(content), os.ModeExclusive)
}

func PrettyPrint(i interface{}) string {
	//s, _ := json.MarshalIndent(i, "", "\t")
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func GetLittleEndianU32Byte(value uint32) []byte {

	key := make([]byte, 4)

	binary.LittleEndian.PutUint32(key[:], value)

	return key
}

func GetLittleEndianU32(value []byte) uint32 {

	return binary.LittleEndian.Uint32(value)
}

func SetLittleEndianU32(buf []byte, value uint32) {
	binary.LittleEndian.PutUint32(buf, value)
}

func SetLittleEndianU16(buf []byte, value uint16) {
	binary.LittleEndian.PutUint16(buf, value)
}

func GenerateRandUnicastMacaddr() string {
	mac := make([]byte, 6)

	// from practical experience, first byte keeps all zero would be the best compatibility
	rand.Read(mac[1:])

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func ContainsInsensitive(haystack, needle string) bool {
	return strings.Contains(
		strings.ToLower(haystack),
		strings.ToLower(needle),
	)
}

func CheckIfSliceContain(slice interface{}, element interface{}) bool {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), element) {
			return true
		}
	}

	return false
}
