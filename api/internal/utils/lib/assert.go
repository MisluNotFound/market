package lib

import (
	"fmt"
	"strings"

	"github.com/mislu/market-api/internal/utils/app"
)

func GetResourceURL(bucketType int, owner string, key string) string {
	return fmt.Sprintf("http://%s/api/assert/%d/%s/%s", app.GetConfig().Server.BaseIP, bucketType, owner, key)
}

func SplitResourceURL(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
