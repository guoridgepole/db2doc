package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd表示调用时没有任何子命令的基本命令
var rootCmd = &cobra.Command{}

// Execute方法将所有子命令添加到根命令并适当设置标志。
// 这是由 main.main()调用的，它只需要在rootCmd上发生一次。
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
