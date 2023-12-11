package misc

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

func Hello() {
	fmt.Println("this is from gokit misc package.")
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
