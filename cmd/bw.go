/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

var ctx, cancel = context.WithCancel(context.Background())

// bwCmd represents the bw command
var bwCmd = &cobra.Command{
	Use:   "bw",
	Short: "Get traffic for network interface",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		netLinks(args[0])
		// 等待接收到中断信号
		waitForInterrupt()
	},
}

func netLinks(networkInterface string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			netStats, err := net.IOCounters(true)
			if err != nil {
				log.Fatal(err)
			}

			// 查找指定网卡的统计信息
			var selectedNetStat *net.IOCountersStat
			for _, stat := range netStats {
				if string(stat.Name) == string(networkInterface) {
					selectedNetStat = &stat
					break
				}

			}
			if selectedNetStat == nil {
				log.Fatalf("网卡 %s 未找到", networkInterface)
			}
			// 打印当前带宽信息（以KB为单位）
			bytesRecv := float64(selectedNetStat.BytesRecv) / 1024.0 / 1024.0
			bytesSent := float64(selectedNetStat.BytesSent) / 1024.0 / 1024.0

			// 打印带宽信息
			fmt.Printf("网卡 %s:\n", networkInterface)
			fmt.Printf("当前接收速率: %.2f MB/s\n", bytesRecv)
			fmt.Printf("当前发送速率: %.2f MB/s\n", bytesSent)

			// 等待一秒再次获取数据
			time.Sleep(time.Second)
		}
	}
}

// waitForInterrupt 等待接收到中断信号
func waitForInterrupt() {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	<-interruptCh

	// 发送取消信号到上下文
	cancel()
}

func init() {
	getCmd.AddCommand(bwCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bwCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bwCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
