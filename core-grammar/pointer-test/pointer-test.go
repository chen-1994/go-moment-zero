package pointer_test

import "fmt"

// 指针
// 取地址符&：一个变量进行取地址，会返回对应类型的指针
// 解引用符*：1.访问指针所指向的元素，也就是解引用	2.声明一个指针
func Test1() {
	num := 2
	p := &num
	fmt.Println(p)
	fmt.Println(*p)

	var numPtr *int     //声明一个指针
	fmt.Println(numPtr) //nil
	numPtr = new(int)   //初始化
	fmt.Println(numPtr)
	fmt.Println(*numPtr)

	fmt.Println(*new(string))    //声明一个字符串指针 初始化 并通过解引符指向元素
	fmt.Println(*new(int))       //声明一个int串指针 初始化 并通过解引符指向元素
	fmt.Println(*new([5]int))    //声明一个数组指针 初始化 并通过解引符指向元素
	fmt.Println(*new([]float64)) //声明一个切片指针 初始化 并通过解引符指向元素
}

// =======================禁止指针运算
func Test2() {
	arr := []int{1, 2, 3, 4, 5}
	p := &arr
	println(&arr[0])
	println(p)
	// 试图进行指针运算
	//p++ //异常
	fmt.Println(p)
}

// new 和 make

//func new(Type) *Type
//返回值是类型指针
//接收参数是类型
//专用于给指针分配内存空间

// func make(t Type, size ...IntegerType) Type
// 返回值是值，不是指针
// 接收的第一个参数是类型，不定长参数根据传入类型的不同而不同
// 专用于给切片，映射表，通道分配内存。
func Test3() {
	a := new(int)                 // int指针
	b := new(string)              // string指针
	c := new([]int)               // 整型切片指针
	d := make([]int, 10, 100)     // 长度为10，容量100的整型切片
	e := make(map[string]int, 10) // 容量为10的映射表
	f := make(chan int, 10)       // 缓冲区大小为10的通道
	fmt.Println(a, b, c, d, e, f)
}
