# checker-middleware

`checker-middleware` 是一个用于检测常见中间件（数据库、缓存、消息队列、对象存储）连通性、写入、删除能力的 CLI 工具。

## 功能特性

- 支持数据库（MySQL、PostgreSQL、达梦、GoldenDB 等）连通性、写入、删除检测
- 支持 Redis（单机、Sentinel、Credis）缓存检测
- 支持 Kafka、RabbitMQ 消息队列检测
- 支持 S3、MinIO、OSS 对象存储检测
- 检查内容包括：连接、写入、删除
- 支持详细 Debug 日志输出

## 安装

```
sh
GOOS=linux GOARCH=amd64  go build -ldflags "-X main.BuildVersion=$(date +%Y%m%d-%H%M%S)" -o checker-middleware main.go
```

## 使用方法

### 查看帮助

```./checker-middleware--help
./checker-middleware --help
```

```
中间件验证CLI工具

Usage:
  checker-middleware [command]

Available Commands:
  cache       验证缓存可用性
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  mq          验证消息队列
  rdb         验证数据库可用性
  storage     验证对象存储

Flags:
  -h, --help      help for checker-middleware
  -v, --version   version for checker-middleware
```

### 数据库可用性检测

```
./checker-middleware rdb --help
```

```
用法:
 checker-middleware rdb [flags]

数据库类型:
 -D, --driver     数据库驱动: [mysql|dm|pgsql|goldendb|mariadb] (default: mysql)

通用参数:
 -d, --db         数据库名
     --debug      Debug模式
 -h, --help       help for rdb
 -H, --host       数据库主机 (default: 127.0.0.1)
 -p, --password   数据库密码
 -P, --port       数据库端口 (default: 3306)
 -u, --user       数据库用户 (default: root)
```

### 缓存可用性检测

```
sh
./checker-middleware cache --help
```

```

用法:
 checker-middleware cache [flags]


缓存类型:
 -m, --mode       Redis模式: [redis|sentinel|credis] (default: redis)

通用参数:
 -d, --db         Redis数据库 (default: 1)
     --debug      Debug模式
 -H, --host       Redis主机 (default: 127.0.0.1)
 -p, --password   Redis密码
 -P, --port       Redis端口 (default: 6379)
 -t, --timeout    连接超时(秒) (default: 10)

Sentinel 专用参数:
 -M, --master     Sentinel主节点名称 (default: mymaster)
 -s, --sentinels  Sentinel主机列表(多主机以,分割) 例: host1:port1,host2:port2
```

### 消息队列可用性检测

```
./checker-middleware mq --help
```

```
用法:
 checker-middleware mq [flags]


消息队列类型:
 -t, --provider   消息队列类型[kafka/rabbitmq] (default: rabbitmq)

通用参数:
     --debug      Debug模式
 -h, --help       help for mq
 -H, --host       RabbitMQ主机
 -p, --password   RabbitMQ密码
 -P, --port       RabbitMQ端口 (default: 5672)
 -u, --user       RabbitMQ用户
 -v, --vhost      RabbitMQ vhost (default: laiye_cloud)

Kafka专用参数:
     --brokers    Kafka地址（host1:port1,host2:port2）
     --topic      Kafka topic (default: laiye_cloud)
```

### 对象存储可用性检测

```
./checker-middleware storage --help
```

```

用法:
 checker-middleware storage [flags]

存储类型:
 -t, --provider   存储类型(s3/oss/minio) (default: minio)

通用参数:
 -u, --access-key AccessKey (default: laiyelaiye)
 -b, --bucket     Bucket
     --debug      Debug
 -H, --endpoint   Endpoint (default: 127.0.0.1:9000)
 -h, --help       help for storage
     --region     Region (default: us-east-1)
 -p, --secret-key SecretKey
     --secure     启用SSL认证
     --timeout    Timeout (default: 10)
     --use-path-style S3请求的URL是否启用路径风格
```

### 参数说明

每个子命令均支持 `--help` 查看详细参数说明。
