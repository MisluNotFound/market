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
	UserBucket BucketType = iota
	ProductBucket
	ConversationBucket
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

func GenerateObjectKey(filename string) string {
	return genUniqueFilename(filename)
}

func GetObjectPath(bucketType BucketType, owner string, key string) string {
	switch bucketType {
	case UserBucket:
		return fmt.Sprintf("users/%s/%s", owner, key)
	case ProductBucket:
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

func UploadFile(bucketType BucketType, path string, data []byte) error {
	switch bucketType {
	default:
		return globalResourceManager.ossClient.Save(path, data)
	}
}

func DeleteFile(BucketType BucketType, key string) error {
	switch BucketType {
	default:
		return globalResourceManager.ossClient.Delete(key)
	}
}

func FileExists(path string) (bool, error) {
	return globalResourceManager.ossClient.Exists(path)
}
