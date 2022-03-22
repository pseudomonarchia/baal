package textcolor

import (
	"fmt"
)

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

// Text color int
const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

// Format formatted output text color
func Format(textColor int, msg string) string {
	conf := 0
	bg := 0

	return fmt.Sprintf(
		"%c[%d;%d;%dm%s%c[0m",
		0x1B,
		conf,
		bg,
		textColor,
		msg,
		0x1B,
	)
}
