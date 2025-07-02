package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/moby/sys/mountinfo"
	"gopkg.in/yaml.v3"
)

var DebugMode bool

func ReadYAML(filename string) (map[string]any, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func TruncateString(s string, max int, enable bool) string {
	// 截断字符串，用...代替
	if !enable {
		return s
	}
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func ToFloat64(v any) float64 {
	// 将任何类型转换为浮点数
	switch num := v.(type) {
	case int:
		return float64(num)
	case int64:
		return float64(num)
	case float64:
		return num
	case string:
		f, _ := strconv.ParseFloat(num, 64)
		return f
	default:
		return 0
	}
}

func IsNumeric(s string) bool {
	// 判断字符串是否是数字
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsVersion(s string) bool {
	// 判断是否是版本号
	if len(s) > 0 && strings.HasPrefix(strings.ToLower(s), "v") {
		s = s[1:]
	}

	return strings.Count(s, ".") >= 1 && strings.IndexFunc(s, func(r rune) bool {
		return !(r >= '0' && r <= '9' || r == '.')
	}) == -1
}

func GetIntPart(parts []string, i int) int {
	// 从字符串切片 parts 中取出第 i 个元素并尝试转为整数，如果越界或转失败返回 0
	if i < len(parts) {
		val, _ := strconv.Atoi(parts[i])
		return val
	}
	return 0
}

func SplitKey(s string, idx int) string {
	// 将字符串切割为类别和项目
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		if idx == 0 {
			return parts[0]
		} else {
			return ""
		}
	}
	if idx == 0 {
		return parts[0] // 类别
	} else {
		return strings.Join(parts[1:], ".") // 项目
	}
}

// 递归获取路径的挂载点
func GetDirMountPoint(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}
	currentPath := absPath
	for currentPath != "/" {
		mount, err := mountinfo.GetMounts(mountinfo.PrefixFilter(currentPath))
		if err != nil {
			return "", fmt.Errorf("failed to get mount info: %v", err)
		}
		if len(mount) > 0 {
			return mount[0].Mountpoint, nil
		}
		currentPath = filepath.Dir(currentPath)
	}
	return "/", nil
}

func TrimProtocol(endpoint string) string {
	// 去除协议头
	protocols := []string{"http://", "https://", "tcp://", "udp://", "s3://"}
	for _, proto := range protocols {
		if after, ok := strings.CutPrefix(endpoint, proto); ok {
			return after
		}
	}
	return endpoint
}
