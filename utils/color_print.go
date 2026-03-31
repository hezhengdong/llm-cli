package utils

import "fmt"

// 固定灰色 ANSI 码（你要的唯一灰色）
const grayColor = "\033[90m"

// 重置颜色码
const resetColor = "\033[0m"

// GrayPrintf 核心函数：用法 = fmt.Printf
// 支持 %s %d %f 所有格式化占位符
func GrayPrintf(format string, args ...interface{}) {
	fmt.Printf(grayColor+"%s"+resetColor, fmt.Sprintf(format, args...))
}
