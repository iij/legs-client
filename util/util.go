package util

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetMacAddr get a mac address from interfaces.
// If device has multiple interfaces, this method returns mac address of first interfaces.
func GetMacAddr() []string {
	ifas, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var macAddrs []string
	for _, ifa := range ifas {
		addr := ifa.HardwareAddr.String()
		if addr != "" {
			macAddrs = append(macAddrs, addr)
		}
	}

	return macAddrs
}

// ToHash create hash from device name and secret.
func ToHash(str string) (hash string) {
	hashBytes := sha512.Sum512([]byte(str))
	hash = hex.EncodeToString(hashBytes[:])

	return
}

// GenID : TODO
func GenID(size int) string {
	// padding すると "=" が出力文字列に含まれる、かつ出力文字列の長さが 4n 固定になってしまうので padding しないようにする。
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	b := make([]byte, enc.DecodedLen(size))
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return enc.EncodeToString(b)
}

// ExpandPath expands a relative path to a absolute path.
func ExpandPath(path string) (string, error) {
	home, err := homedir.Dir()

	if err != nil {
		return "", err
	}
	path = strings.Replace(path, "~", home, 1)
	return filepath.Abs(path)
}
