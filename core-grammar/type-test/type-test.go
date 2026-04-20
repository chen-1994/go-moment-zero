package type_test

import "fmt"

func Test() {
	// ==========静态强类型
	//Go 是一个静态强类型语言，静态指的是 Go 所有变量的类型早在编译期间就已经确定了
	var a int = 60
	fmt.Println(a)

	// ==========类型后置
	//Go 的声明方式始终遵循名字在前面，类型在后面的原则，从左往右读，大概第一眼就可以知道这是一个函数，
	//且返回值为func(int,int) int。当类型变得越来越复杂时，类型后置在可读性上要好得多，Go 在许多层面的设计都是为了可读性而服务的
	//f func(func(int,int) int, int) func(int, int) int

	// ===========类型声明
	//在 Go 中通过类型声明，可以声明一个自定义名称的新类型，声明一个新类型通常需要一个类型名称以及一个基础类型
	type MyInt int64

	type MyFloat64 float64

	type MyMap map[string]int

	// 可以通过编译，但是不建议使用，这会覆盖原有的类型
	type int int64
	//通过类型声明的类型都是新类型，不同的类型无法进行运算，即便基础类型是相同的

	//============类型别名
	type Int = int
	var c Int = 1
	var b int = 2
	//两者是都是同一个类型，仅仅叫的名字不同，所以也就可以进行运算，所以下例自然也就可以通过编译。
	fmt.Println(c + b)
	//内置类型any就是interface{}的类型别名，两者完全等价，仅仅叫法不一样。

	//============类型转换
	//在 Go 中，只存在显式的类型转换，不存在隐式类型转换
	//因此不同类型的变量无法进行运算，无法作为参数传递。类型转换适用的前提是知晓被转换变量的类型和要转换成的目标类型
	var f1 MyFloat64
	var f float64
	f1 = 0.2
	f = 0.1
	//通过显式的将MyFloat64 转换为float64类型，才能进行加法运算。类型转换的另一个前提是：被转换类型必须是可以被目标类型代表的（Representability），
	//例如int可以被int64类型所代表，也可以被float64类型代表，所以它们之间可以进行显式的类型转换，但是int类型无法被string和bool类型代表，因为也就无法进行类型转换。
	fmt.Println(float64(f1) + f)

	var num1 int8 = 1
	var num2 int32 = 512
	//num1被正确的转换为了int32类型，但是num2并没有。这是一个典型的数值溢出问题，int32能够表示 31 位整数，int8只能表示 7 位整数，
	//高精度整数在向低精度整数转换时会抛弃高位保留低位，因此num1的转换结果就是 0。在数字的类型转换中，通常建议小转大，而不建议大转小
	fmt.Println(int32(num1), int8(num2))
	//在使用类型转换时，对于一些类型需要避免歧义，例子如下
	//*Point(p)       // 等价于 *(Point(p))
	//(*Point)(p)     // 将p转换为类型 *Point
	//<-chan int(c)   // 等价于 <-(chan int(c))
	//(<-chan int)(c) // 将c转换为类型  <-chan int
	//(func())(x)     // 将x转换为类型 func()
	//(func() int)(x) // 将x转换为类型 func() int

	// ===============类型断言
	var d int = 1024
	var e interface{} = d
	if val, ok := e.(int); ok {
		fmt.Println(val)
	} else {
		fmt.Println("error type")
	}
	// ===============类型判断
	var g interface{} = 2
	switch g.(type) {
	case int8:
		fmt.Println(g.(int8))
	case int:
		fmt.Println("int")
	case float64:
		fmt.Println("float64")
	case string:
		fmt.Println("string")
	default:
		fmt.Println("unknown")
	}
}

// 这种就没必要使用别名了
type TwoDMap = map[string]map[string]int

func PrintMyMap(mymap map[string]map[string]int) {
	fmt.Println(mymap)
}
func PrintMyMap2(mymap TwoDMap) {
	fmt.Println(mymap)
}
