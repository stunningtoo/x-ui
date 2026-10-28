package service

import (
	"testing"
	"time"
)

// netRate 是 GetStatus 里唯一带除法的计算，历史缺陷是直接 uint64 相减：
// 网卡计数器一旦回退（驱动重载、休眠唤醒后清零），差值下溢成 1.8e19，
// 面板上就是"每秒 16 EB"。这里把异常输入全部钉住。
func TestNetRate(t *testing.T) {
	cases := []struct {
		name     string
		cur      uint64
		last     uint64
		duration time.Duration
		want     uint64
	}{
		{"正常 1 秒 1024 字节", 2048, 1024, time.Second, 1024},
		{"半秒 1024 字节等于 2048 字节每秒", 2048, 1024, 500 * time.Millisecond, 2048},
		{"计数器回退不下溢", 1024, 2048, time.Second, 0},
		{"间隔为零不除零", 2048, 1024, 0, 0},
		{"间隔为负不除负数", 2048, 1024, -time.Second, 0},
		{"计数器无变化", 1024, 1024, time.Second, 0},
		{"异常巨大的差值(超过 1e15)按 0 处理", 1 << 60, 0, time.Second, 0},
	}
	for _, c := range cases {
		if got := netRate(c.cur, c.last, c.duration); got != c.want {
			t.Fatalf("%s: netRate(%d, %d, %v) = %d, 期望 %d", c.name, c.cur, c.last, c.duration, got, c.want)
		}
	}
}
