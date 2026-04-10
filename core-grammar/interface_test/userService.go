package interface_test

import "fmt"

// =============基本接口
type UserInfo interface {
	Get() string
}

type Zhang struct {
	work1 string
}

type Li struct {
	work2 string
}

func (f Zhang) Get() string {
	//fmt.Println("f1")
	return "张三"
}

func (f Li) Get() string {
	//fmt.Println("f2")
	return "李四"
}

type F12 struct {
	UserInfo UserInfo
}

func (c *F12) Build() {
	fmt.Println(c.UserInfo.Get())
}

func Test() {
	f := F12{Zhang{}}
	f.Build()
	fmt.Println()
	f.UserInfo = Li{}
	f.Build()

	z := Zhang{}
	fmt.Println(z.Get())
}

// 特殊例子
// Man是Person超集 所以Man也实现了接口Person，不过这更像是一种"继承"
type Person interface {
	Say(string) string
	Walk(int)
}
type Man interface {
	Exercise()
	Person
}

// 类型Number的底层类型是int，虽然这放在其他语言中看起来很离谱，但Number的方法集确实是Person 的超集，所以也算实现
type Number int

func (n Number) Say(s string) string {
	return "bilibili"
}
func (n Number) Walk(i int) {
	fmt.Println("cat not walk")
}

// 函数类型也可以实现接口
type Func func()

func (f Func) Say(s string) string {
	f()
	return "aaaaaaaa"
}
func (f Func) Walk(i int) {
	f()
	fmt.Println("cat not walk")
}

func Test2() {
	var f Func
	f = func() {
		fmt.Println("hello")
	}
	f()
}

// ========空接口
// 接口内部没有方法集合，根据实现的定义，所有类型都是Any接口的的实现，因为所有类型的方法集都是空集的超集，所以Any接口可以保存任何类型的值
type And interface {
}

func Test3() {
	var a And
	a = 1
	println(a)
	fmt.Println(a)

	a = "aaaa"
	println(a)
	fmt.Println(a)

	a = complex(1, 2)
	println(a)
	fmt.Println(a)

	a = 1.2
	println(a)
	fmt.Println(a)

	a = []int{}
	println(a)
	fmt.Println(a)

	a = map[int]int{}
	println(a)
	fmt.Println(a)
}
