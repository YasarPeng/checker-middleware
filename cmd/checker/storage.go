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

func StorageCheckerCmd() *cobra.Command {
	var (
		storageProvider     string
		storageEndpoint     string
		storageAccess       string
		storageSecret       string
		storageBucket       string
		storageRegion       string
		storageSecure       bool
		storageTimeout      int
		storageDebug        bool
		storageUsePathStyle bool
	)
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "验证对象存储",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug = storageDebug
			cfg := verify.StorageConfig{
				Provider:     storageProvider,
				Endpoint:     storageEndpoint,
				AccessKey:    storageAccess,
				SecretKey:    storageSecret,
				Bucket:       storageBucket,
				Region:       storageRegion,
				Secure:       storageSecure,
				Timeout:      storageTimeout,
				UsePathStyle: storageUsePathStyle,
			}
			result := verify.VerifyStorageJson(cfg)
			fmt.Printf("%s\n", result)
		},
	}
	storageCmd.Flags().StringVarP(&storageProvider, "provider", "t", "minio", "存储类型(s3/oss/minio)")
	storageCmd.Flags().StringVarP(&storageEndpoint, "endpoint", "H", "127.0.0.1:9000", "Endpoint")
	storageCmd.Flags().StringVarP(&storageAccess, "access-key", "u", "laiyelaiye", "AccessKey")
	storageCmd.Flags().StringVarP(&storageSecret, "secret-key", "p", "", "SecretKey")
	storageCmd.Flags().StringVarP(&storageBucket, "bucket", "b", "", "Bucket")
	storageCmd.Flags().StringVar(&storageRegion, "region", "us-east-1", "Region")
	storageCmd.Flags().BoolVar(&storageSecure, "secure", false, "启用SSL认证")
	storageCmd.Flags().IntVar(&storageTimeout, "timeout", 10, "Timeout")
	storageCmd.Flags().BoolVar(&storageDebug, "debug", false, "Debug")
	storageCmd.Flags().BoolVar(&storageUsePathStyle, "use-path-style", false, "S3请求的URL是否启用路径风格")

	storageCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("\n用法:")
		fmt.Printf(" %s\n", cmd.UseLine())
		fmt.Println("\n存储类型:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"provider"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\n通用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"provider"}
			if !slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
	})
	return storageCmd
}
