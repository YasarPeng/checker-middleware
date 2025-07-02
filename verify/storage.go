package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"checker-middleware/pkg/logger"
	pkgutil "checker-middleware/pkg/util"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// 连通性测试
func StorageConnect(cfg StorageConfig) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "minio":
		client, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.Secure,
		})
		if err != nil {
			result["error"] = fmt.Sprintf("minio connect error: %v", err)
			logger.DebugLog("minio connect error: %v", err)
			return result
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		exists, err := client.BucketExists(ctx, cfg.Bucket)
		if err != nil || !exists {
			result["error"] = fmt.Sprintln("Bucket not exists")
			logger.DebugLog("Bucket not exists")
			return result
		}
		result["success"] = "true"
		logger.DebugLog("minio connect success")
	case "oss":
		client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			result["error"] = fmt.Sprintf("oss connect error: %v", err)
			logger.DebugLog("oss connect error: %v", err)
			return result
		}
		bucket, err := client.Bucket(cfg.Bucket)
		if err != nil {
			result["error"] = fmt.Sprintf("oss bucket error: %v", err)
			logger.DebugLog("oss bucket error: %v", err)
			return result
		}
		_, _ = bucket.GetObjectMeta("not_exist_key")
		// 只要能正常请求即可
		result["success"] = "true"
		logger.DebugLog("oss connect success")
	case "s3":
		endpoint := pkgutil.TrimProtocol(cfg.Endpoint)
		protocol := "http://"
		if cfg.Secure {
			protocol = "https://"
		}
		awsCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(awsCreds.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithEndpointResolver(
				aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               protocol + endpoint,
						SigningRegion:     cfg.Region,
						HostnameImmutable: true,
					}, nil
				}),
			),
		)
		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 config error: %v", err)
			logger.DebugLog("s3 v2 config error: %v", err)
			return result
		}
		client := awsS3.NewFromConfig(awsCfg, func(o *awsS3.Options) {
			o.UsePathStyle = cfg.UsePathStyle
		})
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.Timeout)*time.Second,
		)
		defer cancel()
		_, err = client.HeadBucket(ctx, &awsS3.HeadBucketInput{
			Bucket: &cfg.Bucket,
		})
		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 bucket error: %v", err)
			logger.DebugLog("s3 v2 bucket error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("s3 v2 connect success")
	// ...existing code...
	default:
		result["error"] = "unsupported provider"
	}
	return result
}

// 上传测试
func StorageWrite(cfg StorageConfig, objectName, content string) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "minio":
		client, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.Secure,
		})
		if err != nil {
			result["error"] = fmt.Sprintf("minio connect error: %v", err)
			logger.DebugLog("minio connect error: %v", err)
			return result
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		_, err = client.PutObject(ctx, cfg.Bucket, objectName, bytes.NewReader([]byte(content)), int64(len(content)), minio.PutObjectOptions{})
		if err != nil {
			result["error"] = fmt.Sprintf("minio upload error: %v", err)
			logger.DebugLog("minio upload error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("minio upload success")
	case "oss":
		client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			result["error"] = fmt.Sprintf("oss connect error: %v", err)
			logger.DebugLog("oss connect error: %v", err)
			return result
		}
		bucket, err := client.Bucket(cfg.Bucket)
		if err != nil {
			result["error"] = fmt.Sprintf("oss bucket error: %v", err)
			logger.DebugLog("oss bucket error: %v", err)
			return result
		}
		err = bucket.PutObject(objectName, bytes.NewReader([]byte(content)))
		if err != nil {
			result["error"] = fmt.Sprintf("oss upload error: %v", err)
			logger.DebugLog("oss upload error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("oss upload success")
	case "s3":
		endpoint := pkgutil.TrimProtocol(cfg.Endpoint)
		protocol := "http://"
		if cfg.Secure {
			protocol = "https://"
		}
		awsCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(awsCreds.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithEndpointResolver(
				aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               protocol + endpoint,
						SigningRegion:     cfg.Region,
						HostnameImmutable: true,
					}, nil
				}),
			),
		)
		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 config error: %v", err)
			logger.DebugLog("s3 v2 config error: %v", err)
			return result
		}
		client := awsS3.NewFromConfig(awsCfg, func(o *awsS3.Options) {
			o.UsePathStyle = cfg.UsePathStyle
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		_, err = client.PutObject(ctx, &awsS3.PutObjectInput{
			Bucket: &cfg.Bucket,
			Key:    &objectName,
			Body:   bytes.NewReader([]byte(content)),
		})
		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 upload error: %v", err)
			logger.DebugLog("s3 v2 upload error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("s3 v2 upload success")
	// ...existing code...
	default:
		result["error"] = "unsupported provider"
	}
	return result
}

// 删除测试
func StorageDelete(cfg StorageConfig, objectName string) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "minio":
		client, err := minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.Secure,
		})
		if err != nil {
			result["error"] = fmt.Sprintf("minio connect error: %v", err)
			logger.DebugLog("minio connect error: %v", err)
			return result
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		err = client.RemoveObject(ctx, cfg.Bucket, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			result["error"] = fmt.Sprintf("minio delete error: %v", err)
			logger.DebugLog("minio delete error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("minio delete success")
	case "oss":
		client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			result["error"] = fmt.Sprintf("oss connect error: %v", err)
			logger.DebugLog("oss connect error: %v", err)
			return result
		}
		bucket, err := client.Bucket(cfg.Bucket)
		if err != nil {
			result["error"] = fmt.Sprintf("oss bucket error: %v", err)
			logger.DebugLog("oss bucket error: %v", err)
			return result
		}
		err = bucket.DeleteObject(objectName)
		if err != nil {
			result["error"] = fmt.Sprintf("oss delete error: %v", err)
			logger.DebugLog("oss delete error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("oss delete success")
	case "s3":
		endpoint := pkgutil.TrimProtocol(cfg.Endpoint)
		protocol := "http://"
		if cfg.Secure {
			protocol = "https://"
		}
		awsCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(awsCreds.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithEndpointResolver(
				aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               protocol + endpoint,
						SigningRegion:     cfg.Region,
						HostnameImmutable: true,
					}, nil
				}),
			),
		)

		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 config error: %v", err)
			logger.DebugLog("s3 v2 config error: %v", err)
			return result
		}
		client := awsS3.NewFromConfig(awsCfg, func(o *awsS3.Options) {
			o.UsePathStyle = cfg.UsePathStyle
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		_, err = client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
			Bucket: &cfg.Bucket,
			Key:    &objectName,
		})
		if err != nil {
			result["error"] = fmt.Sprintf("s3 v2 delete error: %v", err)
			logger.DebugLog("s3 v2 delete error: %v", err)
			return result
		}
		result["success"] = "true"
		logger.DebugLog("s3 v2 delete success")
	// ...existing code...
	default:
		result["error"] = "unsupported provider"
	}
	return result
}

// 一键检测
func VerifyStorage(cfg StorageConfig) StorageResult {
	objectName := "precheck_test_object.txt"
	content := "hello precheck"
	res := StorageResult{
		Connect: StorageConnect(cfg),
		Write:   map[string]string{"success": "skip"},
		Delete:  map[string]string{"success": "skip"},
	}
	res.Write = StorageWrite(cfg, objectName, content)
	if res.Connect["success"] == "true" {
		res.Write = StorageWrite(cfg, objectName, content)
		if res.Write["success"] == "true" {
			res.Delete = StorageDelete(cfg, objectName)
		}
	}
	return res
}

func VerifyStorageJson(cfg StorageConfig) []byte {
	res := VerifyStorage(cfg)
	r, _ := json.Marshal(res)
	return r
}
