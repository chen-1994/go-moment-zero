package func_test

import (
	"fmt"
	"io"
	"slices"
)

// =============声明
// 函数的两种 声明方法 func var
func sum(a int, b int) int {
	return a + b
}

// go 不支持重载 函数名称相同不行
var sum1 = func(a int, b int) int {
	return a + b
}

// ===========参数
// go参数可以不带名称，一般在接口或函数类型声明中
type ExWriter func(writer io.Writer) error

type Writer interface {
	Write([]byte) (int, error)
}

// 变长参数可以接收0或者多个值，只能在参数列表末尾
func Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}

// ==============匿名函数
func Test1() {
	func() int {
		return 1
	}()
	func(a int) int {
		return a
	}(1)
}

type Person struct {
	Name   string
	Age    int
	Salary float64
}

func Test2() {
	people := []Person{
		{Name: "Alice", Age: 25, Salary: 5000.0},
		{Name: "Bob", Age: 30, Salary: 6000.0},
		{Name: "Charlie", Age: 28, Salary: 5500.0},
	}
	// 排序
	slices.SortFunc(people, func(p1 Person, p2 Person) int {
		if p1.Age > p2.Age {
			return 1
		} else if p1.Age < p2.Age {
			return -1
		}
		return 0
	})
	fmt.Println(people)
}

// =========闭包
// Exp(2) 本质上创建了一个函数grow 初始花了e 函数内容为匿名函数 return的返回值 循环10次，匿名函数被调用了10次
func Test3() {
	grow := Exp(2)
	for i := range 10 {
		fmt.Printf("2^%d=%d\n\n", i, grow())
	}
}
func Exp(n int) func() int {
	e := 1
	return func() int {
		tmpe := e
		e *= n

		return tmpe
	}
}

// ======延迟调用 延迟都是最后执行，多个延迟 先进后出
func Test4() {
	//defer1()
	//fmt.Println(defer2(1, 3))
	//defer3()
	defer4()
}
func defer1() {
	defer func() {
		fmt.Println("defer1 1")
	}()
	fmt.Println("defer1 2")
}

func defer2(a, b int) (s int) {
	defer func() {
		s -= 10
	}()
	s = a + b
	return
}
func defer3() {
	defer fmt.Println("defer3 1")
	fmt.Println("defer3 2")
	defer fmt.Println("defer3 3")
	defer fmt.Println("defer3 4")
	fmt.Println("defer3 5")
}
func defer4() {
	for i := range 5 {
		defer fmt.Println(i)
	}
}

// =======参数预设 defer 函数不会执行，但参数会先执行，如果参数是函数，那也会先执行
func Test5() {
	//Fn()
	//s5()
	//s6()
	//s7()
	s8()
}
func Fn() {
	defer fmt.Println(Fn1())
	fmt.Println("3")
}
func Fn1() int {
	fmt.Println("2")
	return 1
}

// 参数改变
func s5() {
	var a, b int = 1, 2
	defer fmt.Println(sum5(a, b))
	a, b = 3, 4
}
func sum5(a int, b int) int {
	return a + b
}

// 闭包 f := func() 只是创建了一个名为f的函数 f()是函数正式执行，参数a, b = 3, 4
func s6() {
	var a, b = 1, 2
	f := func() {
		fmt.Println(sum5(a, b))
	}
	a, b = 3, 4
	f()
}

// 延迟与闭包结合
// func 匿名函数sum5(a, b)不是参数，所以不会触发defer中的函数的参数先执行，defer等待主函数执行完再执行 理解正确吗
func s7() {
	var a, b = 1, 2
	defer func() {
		fmt.Println(sum5(a, b))
	}()
	a, b = 3, 4
}

// sum5(a, b) 作为参数传递，会先执行
func s8() {
	var a, b = 1, 2
	defer func(num int) {
		fmt.Println(num)
	}(sum5(a, b))
	a, b = 3, 4
}
