package converter

import "strings"

func BuildAssetURL(baseURL, objectKey string) string {
	key := strings.TrimPrefix(objectKey, "/")
	return strings.TrimRight(baseURL, "/") + "/" + key
}
