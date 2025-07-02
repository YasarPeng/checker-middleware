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

func CacheCheckerCMD() *cobra.Command {

	var (
		cacheHost      string
		cachePort      int
		cachePassword  string
		cacheDebug     bool
		cacheDB        int
		cacheMode      string
		cacheSentinels []string
		cacheMaster    string
		cacheTimeout   int
	)
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "验证缓存可用性",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := verify.CacheConfig{
				Host:      cacheHost,
				Port:      cachePort,
				Password:  cachePassword,
				DB:        cacheDB,
				Mode:      cacheMode,
				Sentinels: cacheSentinels,
				Master:    cacheMaster,
				Timeout:   cacheTimeout,
			}
			logger.Debug = cacheDebug
			result := verify.VerifyCacheJson(cfg)
			fmt.Println(string(result))
		},
	}
	// 通用参数
	cacheCmd.Flags().StringVarP(&cacheHost, "host", "H", "127.0.0.1", "Redis主机")
	cacheCmd.Flags().IntVarP(&cachePort, "port", "P", 6379, "Redis端口")
	cacheCmd.Flags().StringVarP(&cachePassword, "password", "p", "", "Redis密码")
	cacheCmd.Flags().IntVarP(&cacheDB, "db", "d", 1, "Redis数据库")
	cacheCmd.Flags().BoolVar(&cacheDebug, "debug", false, "Debug模式")
	cacheCmd.Flags().IntVarP(&cacheTimeout, "timeout", "t", 10, "连接超时(秒)")
	cacheCmd.Flags().StringVarP(&cacheMode, "mode", "m", "redis", "Redis模式: [redis|sentinel|credis]")

	// Sentinel专用参数
	cacheCmd.Flags().StringSliceVarP(&cacheSentinels, "sentinels", "s", []string{}, "Sentinel主机列表(多主机以,分割) 例: host1:port1,host2:port2")
	cacheCmd.Flags().StringVarP(&cacheMaster, "master", "M", "mymaster", "Sentinel主节点名称")
	// 自定义帮助信息
	cacheCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("\n用法:")
		fmt.Printf(" %s\n\n", cmd.UseLine())
		fmt.Println("\n缓存类型:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"mode"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\n通用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"host", "port", "password", "db", "timeout", "debug"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\nSentinel 专用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"sentinels", "master"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
	})

	return cacheCmd
}
