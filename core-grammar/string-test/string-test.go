package string_test

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"unsafe"
)

func StringTest1() {
	a := "1 2 33"
	b := `
			"你好啊 朋友"
		`
	fmt.Println(a, b)
	fmt.Println(a[0])         //返回编码值
	fmt.Println(string(a[0])) //string 可以转原值
	//a[0] = 'a' // 无法通过编译

	bytes := []byte(a)
	bytes = append(bytes, 96, 97, 98) //97,98,99 是编码值
	fmt.Println(a)                    //字符串切片 不会影响原来的字符串。切片为新数据内存 内存安全
	a = string(bytes)
	fmt.Println(a)
}

func StringTest2() {

	//unsafe 提供无复制转换
	a := "hello world"
	b := unsafe.Slice(unsafe.StringData(a), len(a))
	fmt.Println("%p %p", unsafe.StringData(a), unsafe.StringData(string(b)))

	//字符串copy
	c := strings.Clone(a) //copy 后的地址不一样
	fmt.Println(c, unsafe.StringData(c))

	builder := strings.Builder{}
	builder.WriteString(a)
	builder.WriteString("你好啊")
	fmt.Println(builder.String())

	str := builder.String()
	for i := 0; i < len(str); i++ { //中文会乱码
		fmt.Printf("%d,%x,%s\n", str[i], str[i], string(str[i]))
	}
	//可以使用 range
	for _, v := range str {
		fmt.Printf("%d,%x,%s\n", v, v, string(v))
	}
	//可以使用rune
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		fmt.Println(string(runes[i]))
	}
	//使用utf8工具包
	for i, w := 0, 0; i < len(str); i += w {
		r, width := utf8.DecodeRuneInString(str[i:])
		fmt.Println(string(r))
		w = width
	}
}
