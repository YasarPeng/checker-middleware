package checker

import (
	"fmt"
	"checker-middleware/pkg/logger"
	pkgutil "checker-middleware/pkg/util"
	"checker-middleware/verify"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func MQCheckerCmd() *cobra.Command {
	// verify mq
	var (
		mqProvider string
		mqBrokers  string
		mqTopic    string
		mqHost     string
		mqPort     int
		mqUser     string
		mqPassword string
		mqVhost    string
		mqDebug    bool
	)
	mqCmd := &cobra.Command{
		Use:   "mq",
		Short: "验证消息队列",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := verify.MQConfig{
				Provider: mqProvider,
				Brokers:  strings.Split(mqBrokers, ","),
				Topic:    mqTopic,
				Host:     mqHost,
				Port:     mqPort,
				User:     mqUser,
				Password: mqPassword,
				Vhost:    mqVhost,
			}
			logger.Debug = mqDebug
			result := verify.VerifyMQJson(cfg)
			fmt.Println(result)
		},
	}
	mqCmd.Flags().StringVarP(&mqProvider, "provider", "t", "rabbitmq", "消息队列类型[kafka/rabbitmq]")
	mqCmd.Flags().StringVar(&mqBrokers, "brokers", "", "Kafka地址（host1:port1,host2:port2）")
	mqCmd.Flags().StringVar(&mqTopic, "topic", "laiye_cloud", "Kafka topic")
	mqCmd.Flags().StringVarP(&mqHost, "host", "H", "", "RabbitMQ主机")
	mqCmd.Flags().IntVarP(&mqPort, "port", "P", 5672, "RabbitMQ端口")
	mqCmd.Flags().StringVarP(&mqUser, "user", "u", "", "RabbitMQ用户")
	mqCmd.Flags().StringVarP(&mqPassword, "password", "p", "", "RabbitMQ密码")
	mqCmd.Flags().StringVarP(&mqVhost, "vhost", "v", "laiye_cloud", "RabbitMQ vhost")
	mqCmd.Flags().BoolVar(&mqDebug, "debug", false, "Debug模式")

	mqCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("\n用法:")
		fmt.Printf(" %s\n\n", cmd.UseLine())
		fmt.Println("\n消息队列类型:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"provider"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\n通用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"brokers", "topic", "provider"}
			if !slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
		fmt.Println("\nKafka专用参数:")
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			names := []string{"brokers", "topic"}
			if slices.Contains(names, f.Name) {
				pkgutil.PrintFlag(f)
			}
		})
	})

	return mqCmd
}
