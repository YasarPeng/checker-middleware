package checker

import (
	"fmt"
	"checker-middleware/pkg/logger"
	pkgutil "checker-middleware/pkg/util"
	"checker-middleware/verify"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func RDBCheckerCmd() *cobra.Command {
	var (
		rdbDriver   string
		rdbHost     string
		rdbPort     int
		rdbUser     string
		rdbPassword string
		rdbDatabase string
		rdbDebug    bool
	)
	rdbCmd := &cobra.Command{
		Use:   "rdb",
		Short: "验证数据库可用性",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := verify.RDBConfig{
				Driver:   rdbDriver,
				Host:     rdbHost,
				Port:     rdbPort,
				Username: rdbUser,
				Password: rdbPassword,
				Database: rdbDatabase,
			}
			logger.Debug = rdbDebug
			result := verify.VerifyRDBJson(cfg)
			fmt.Printf("VerifyRDBJson: %s\n", result)
		},
	}
	rdbCmd.Flags().StringVarP(&rdbDriver, "driver", "D", "mysql", "数据库驱动: [mysql|dm|pgsql|goldendb|mariadb]")
	rdbCmd.Flags().StringVarP(&rdbHost, "host", "H", "127.0.0.1", "数据库主机")
	rdbCmd.Flags().IntVarP(&rdbPort, "port", "P", 3306, "数据库端口")
	rdbCmd.Flags().StringVarP(&rdbUser, "user", "u", "root", "数据库用户")
	rdbCmd.Flags().StringVarP(&rdbPassword, "password", "p", "", "数据库密码")
	rdbCmd.Flags().StringVarP(&rdbDatabase, "db", "d", "", "数据库名")
	rdbCmd.Flags().BoolVar(&rdbDebug, "debug", false, "Debug模式")

	rdbCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("\n用法:")
		fmt.Printf(" %s\n", cmd.UseLine())
		fmt.Println("\n数据库类型:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"driver"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\n通用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"driver"}
			if !slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
	})
	return rdbCmd
}
