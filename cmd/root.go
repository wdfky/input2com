package cmd

import (
	"fmt"
	"input2com/internal/cli"
	"os"

	"github.com/spf13/cobra"
)

import (
	"input2com/internal/config"
)

// RootCmd 代表根命令
var RootCmd = &cobra.Command{
	Use:   "input2com",
	Short: "将输入设备事件转发到串口",
	Long:  `一个用于将鼠标、键盘、手柄等输入设备事件通过串口转发出去的工具。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.Run(config.Cfg.Debug, config.Cfg.Baudrate, config.Cfg.TtyPath, config.Cfg.MouseConfigDict)
		return nil
	},
}

// Execute 由 main.go 调用
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
}
