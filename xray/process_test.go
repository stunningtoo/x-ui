package xray

import (
	"errors"
	"testing"
	"time"
)

// GetResult 在空日志队列上必须立即返回。
// 历史实现无条件走 TakeUntil，空队列会永久阻塞；调用方 GetXrayResult
// 持有 xrayLock 写锁，卡住之后每 10 秒的 RestartXray 和 /server/status
// 会一起被拖死，面板表现为 xray 状态永远转圈、点重启没反应。
func TestGetResultEmptyQueueNotBlocking(t *testing.T) {
	p := newProcess(&Config{})
	done := make(chan string, 1)
	go func() {
		done <- p.GetResult()
	}()
	select {
	case got := <-done:
		if got != "" {
			t.Fatalf("空队列且无错误时应返回空串，实际: %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("GetResult 在空队列上阻塞了")
	}
}

// 空队列但有退出错误时，返回错误文本而不是阻塞。
func TestGetResultEmptyQueueReturnsExitErr(t *testing.T) {
	p := newProcess(&Config{})
	p.exitErr = errors.New("exit status 1")

	done := make(chan string, 1)
	go func() {
		done <- p.GetResult()
	}()
	select {
	case got := <-done:
		if got != "exit status 1" {
			t.Fatalf("应返回退出错误文本，实际: %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("GetResult 在空队列上阻塞了")
	}
}

// 有日志时按顺序拼接全部行，且取完即清空。
func TestGetResultReturnsQueuedLines(t *testing.T) {
	p := newProcess(&Config{})
	p.lines.Put("started")
	p.lines.Put("listening on 127.0.0.1:62789")

	got := p.GetResult()
	want := "started\nlistening on 127.0.0.1:62789"
	if got != want {
		t.Fatalf("GetResult() = %q, 期望 %q", got, want)
	}
	if !p.lines.Empty() {
		t.Fatal("取过的日志应从队列中移除")
	}
}
