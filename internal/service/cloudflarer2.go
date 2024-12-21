package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"openapi/internal/constants"
	"openapi/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// createR2Client 创建 R2 客户端
func createR2Client() (*s3.Client, error) {
	// 获取配置（使用缓存）
	config, err := GetCloudflareConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get cloudflare config: %v", err)
	}

	// 创建 AWS 配置
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.AccessKeySecret,
			"",
		)),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	// 创建 S3 客户端
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com",
			config.AccountID))
	})

	return client, nil
}

// GetBucketInfo 获取 bucket 目录信息
func GetBucketInfo(countryCode, env, bucketName string) (*constants.CFResponse[constants.KVKeys], error) {
	client, err := createR2Client()
	if err != nil {
		return nil, err
	}

	// 列出对象，使用分隔符获取目录
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Delimiter: aws.String("/"), // 使用 / 作为分隔符来获取目录
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	// 构建响应
	var result constants.CFResponse[constants.KVKeys]
	result.Success = true
	result.Result = make([]constants.KVKeys, 0)

	// 添加公共前缀（目录）信息
	for _, prefix := range listObjectsOutput.CommonPrefixes {
		// 移除末尾的 /
		dirName := strings.TrimSuffix(*prefix.Prefix, "/")
		result.Result = append(result.Result, constants.KVKeys{
			Name: dirName,
		})
	}

	// 添加根目录下的文件（可选，如果只需要目录可以注释掉）
	for _, object := range listObjectsOutput.Contents {
		// 跳过包含 / 的对象（子目录中的文件）
		if !strings.Contains(*object.Key, "/") {
			result.Result = append(result.Result, constants.KVKeys{
				Name: *object.Key,
			})
		}
	}

	// 打印调试信息
	logger.Info("Found %d directories and root files in bucket %s",
		len(result.Result), bucketName)

	return &result, nil
}

// CopyTask 表示一个复制任务
type CopyTask struct {
	SourceKey string
	TargetKey string
}

// CopyDirectory 复制目录及其内容到新位置
func CopyDirectory(countryCode, Env, bucketName, sourceDir, targetDir string) error {
	client, err := createR2Client()
	if err != nil {
		return err
	}

	// 确保源目录和目标目录以 / 结尾
	sourceDir = strings.TrimSuffix(sourceDir, "/") + "/"
	targetDir = strings.TrimSuffix(targetDir, "/") + "/"

	// 创建任务通道和错误通道
	tasks := make(chan CopyTask, 100)
	errors := make(chan error, 1)
	var wg sync.WaitGroup

	// 启动工作池
	workerCount := 50 // 并发工作器数量
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				// 创建复制源引用
				copySource := fmt.Sprintf("%s/%s", bucketName, task.SourceKey)

				// 复制对象
				_, err := client.CopyObject(context.TODO(), &s3.CopyObjectInput{
					Bucket:     aws.String(bucketName),
					CopySource: aws.String(copySource),
					Key:        aws.String(task.TargetKey),
				})
				if err != nil {
					select {
					case errors <- fmt.Errorf("failed to copy %s to %s: %v",
						task.SourceKey, task.TargetKey, err):
					default:
					}
					return
				}

				logger.Info("Copied %s to %s", task.SourceKey, task.TargetKey)
			}
		}()
	}

	// 列出源目录中的所有对象并创建任务
	go func() {
		defer close(tasks)
		paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String(sourceDir),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				errors <- fmt.Errorf("failed to list objects: %v", err)
				return
			}

			for _, obj := range page.Contents {
				sourceKey := *obj.Key
				targetKey := targetDir + strings.TrimPrefix(sourceKey, sourceDir)
				tasks <- CopyTask{
					SourceKey: sourceKey,
					TargetKey: targetKey,
				}
			}
		}
	}()

	// 等待所有工作完成
	wg.Wait()

	// 检查是否有错误
	select {
	case err := <-errors:
		return err
	default:
	}

	logger.Info("Successfully copied directory from %s to %s", sourceDir, targetDir)
	return nil
}

// DeleteTask 表示一个删除任务
type DeleteTask struct {
	Key string
}

// DeleteDirectory 删除目录及其内容
func DeleteDirectory(countryCode, Env, bucketName, dirPath string) error {
	client, err := createR2Client()
	if err != nil {
		return err
	}

	// 确保目录路径以 / 结尾
	dirPath = strings.TrimSuffix(dirPath, "/") + "/"

	// 创建任务通道和错误通道
	tasks := make(chan DeleteTask, 100)
	errors := make(chan error, 1)
	var wg sync.WaitGroup

	// 启动工作池
	workerCount := 50 // 并发工作器数量
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				// 删除对象
				_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(task.Key),
				})
				if err != nil {
					select {
					case errors <- fmt.Errorf("failed to delete object %s: %v", task.Key, err):
					default:
					}
					return
				}

				logger.Info("Deleted %s", task.Key)
			}
		}()
	}

	// 列出目录中的所有对象并创建删除任务
	go func() {
		defer close(tasks)
		paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String(dirPath),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				errors <- fmt.Errorf("failed to list objects: %v", err)
				return
			}

			for _, obj := range page.Contents {
				tasks <- DeleteTask{
					Key: *obj.Key,
				}
			}
		}
	}()

	// 等待所有工作完成
	wg.Wait()

	// 检查是否有错误
	select {
	case err := <-errors:
		return err
	default:
	}

	logger.Info("Successfully deleted directory %s", dirPath)
	return nil
}
