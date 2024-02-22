package cmd

import (
	"bytes"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetLinks(t *testing.T) {
	var buf bytes.Buffer
	// 保存原始的标准输出
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// 使用测试的网络接口名称
	networkInterface := "WLAN"

	// 启动测试，等待一秒钟以获取一些输出
	go netLinks("WLAN")
	time.Sleep(time.Second)

	// 结束测试，发送中断信号
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT)
	interruptCh <- syscall.SIGINT

	// 等待一些时间确保程序有足够的时间响应中断信号
	time.Sleep(time.Second * 5)

	// 检查输出是否包含了预期的信息
	// 这里使用你的网络接口名称和预期的输出来替换以下断言中的占位符
	output := buf.String()
	assert.Contains(t, output, "网卡 "+networkInterface)
	assert.Contains(t, output, "当前接收速率:")
	assert.Contains(t, output, "当前发送速率:")
}
