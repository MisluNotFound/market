package resourcemanager

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"time"

	"github.com/mislu/market-api/internal/oss"
	"github.com/mislu/market-api/internal/oss/local"
	"github.com/mislu/market-api/internal/utils/app"
)

type BucketType int

const (
	UserAvatarBucket BucketType = iota
	ProductImageBucket
	ProductVideoBucket
	ProductAvatarBucket 
	ChatSessionBucket
	TempUploadBucket
)

var globalResourceManager ResourceManager

const (
	LOCAL_STORAGE = "local_storage"
	// TODO xx oss
)

type ResourceManager struct {
	ossClient oss.OSS

	// securityCli   ContentSecurity
	// cdnDomains    map[BucketType]string cdns
	// imgProcessor  ImageProcessor  process image to different size
	// costCollector CostCollector   content analysis
}

func InitGlobalResourceManager() {
	config := app.GetConfig().OSS
	
	var ossClient oss.OSS
	switch config.Type {
	case LOCAL_STORAGE:
		ossClient = local.NewLocalStorage(config.Root)
	}

	globalResourceManager.ossClient = ossClient
}

func GenerateObjectKey(bucketType BucketType, owner string, filename string) string {
	return GetObjectPath(bucketType, owner, genUniqueFilename(filename))
}

func GetObjectPath(bucketType BucketType, owner string, key string) string {
	switch bucketType {
	case UserAvatarBucket:
		return fmt.Sprintf("users/%s/%s", owner, key)
	case ProductImageBucket:
		return fmt.Sprintf("products/%s/%s", owner, key)
	case ProductVideoBucket:
		return fmt.Sprintf("products/%s/%s", owner, key)
	case ProductAvatarBucket:
		return fmt.Sprintf("products/%s/%s", owner, key)
	default:
		return fmt.Sprintf("temp/%s", key)
	}
}


func genUniqueFilename(origin string) string {
	ext := filepath.Ext(origin)
	hash := sha256.Sum256([]byte(origin + time.Now().String()))
	return fmt.Sprintf("%x%s", hash[:8], ext)
}

func UploadFile(bucketType BucketType, key string, data []byte) error {
	// TODO 在bucket内部检查文件大小
	switch bucketType {
	case UserAvatarBucket:
		return globalResourceManager.ossClient.Save(key, data)
	}

	return nil
}

func DeleteFile(BucketType BucketType, key string) error {
	switch BucketType {
	case UserAvatarBucket:
		return globalResourceManager.ossClient.Delete(key)
	}

	return nil
}