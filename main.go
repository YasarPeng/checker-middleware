package main

import (
	checker "checker-middleware/cmd/checker"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	BuildVersion = "20250701-1828"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "checker-middleware",
		Short: "中间件验证CLI工具",
	}
	rdbCmd := checker.RDBCheckerCmd()
	cacheCmd := checker.CacheCheckerCMD()
	mqCmd := checker.MQCheckerCmd()
	storageCmd := checker.StorageCheckerCmd()

	rootCmd.Version = BuildVersion
	rootCmd.AddCommand(rdbCmd, cacheCmd, storageCmd, mqCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
