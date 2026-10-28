//go:build windows
// +build windows

package sys

import (
	"github.com/shirou/gopsutil/net"
)

// Windows 没有 /proc/net，只能走 gopsutil 的 iphlpapi 实现。
// 返回整机所有 TCP/UDP 连接数，语义与 linux/darwin 版本一致。
func GetTCPCount() (int, error) {
	stats, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}
	return len(stats), nil
}

func GetUDPCount() (int, error) {
	stats, err := net.Connections("udp")
	if err != nil {
		return 0, err
	}
	return len(stats), nil
}
