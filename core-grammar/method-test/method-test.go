package method_test

import "fmt"

type IntSlice []int

func (p IntSlice) Get(index int) int {
	return p[index]
}
func (p IntSlice) Set(index int, value int) {
	p[index] = value
}

func (P IntSlice) len() int {
	return len(P)
}

func Test1() {
	p := IntSlice{1, 2, 3}
	fmt.Println(p.Get(1))
	p.Set(2, 1)
	fmt.Println(p)
	fmt.Println(p.len())
}

// 值接收者 值接收者无法变更原值
type MyInt int

func (p MyInt) Set(val int) {
	p = MyInt(val)
}
func Test2() {
	myInt := MyInt(1)
	myInt.Set(10)
	fmt.Println(myInt)

	(&myInt).Set(10)
	fmt.Println(myInt)
}

// 指针接收者
// 函数的参数传递的过程中 是值拷贝的。如果传递的是一个整型，那就拷贝整型。如果是一个切片就拷贝一个切片。但如果是一个指针就只需要拷贝这个指针。拷贝
// 指针比拷贝切片消耗的资源更小，接收者也不例外。值接收者和指针接收者也是同样的道理。大多数情况下推荐使用指针接收者。不过两者不应该混用。要么都用指针，要么都不用指针
func (p *MyInt) Set2(val int) {
	*p = MyInt(val)
}
func Test3() {
	myInt := MyInt(1)
	myInt.Set2(10)
	fmt.Println(myInt)
}

type Animal interface {
	Run()
	Run2()
}
type Dog struct {
}

func (d *Dog) Run() {
	fmt.Println("Run")
}
func (d Dog) Run2() {
	fmt.Println("Run2")
}
func Test4() {
	var a Animal
	a = &Dog{}
	//a = Dog{} 异常 因为接口实现中有指针引用
	a.Run()
	a.Run2()
}

// 当值接收者可以寻址，Go 会自动的插入指针运算符来进行调用，例如切片是可寻址，依旧可以通过值接收者来修改其内部值
//
//	语法糖：自动寻址 (Automatic Pointer Conversion)
//
// “自动插入指针运算符”在 Go 中确实存在，但这不是针对特定类型（如切片或 Map），而是针对变量是否可寻址。
// 规则：如果你有一个变量 v，即使它不是指针，只要它是可寻址的（Addressable），你就可以直接调用指针接收者的方法 v.PointerMethod()。Go 会自动帮你转成 (&v).PointerMethod()。
// 哪些不可寻址？：
// Map 的元素：myMap["key"].Set() 会报错（因为 Map 扩容会导致元素地址变动）。
// 常量、字面量：MyInt(10).Set() 会报错。

type Slice []int

func (s Slice) Set(i, v int) {
	s[i] = v
}
func Test5() {
	s := Slice{1, 2, 3}
	s.Set(1, 10)
	fmt.Println(s)
}

// 坑点A：Slice 的 append 问题：只要涉及到改变长度（扩容），必须用指针接收者
// 坑 B：Map 的元素不可寻址 原因：Map 内部为了性能会频繁重排内存，导致元素的地址不固定。所以 Go 禁止直接取 Map 元素的地址。
// 解决：Map 里面存指针 map[int]*User，或者取出整个结构体修改后再塞回去
func (s Slice) Add(v int) {
	s = append(s, v) // 这里改的是拷贝出来的 len/cap，原切片长度不变
}
func Test6() {
	s := Slice{1, 2, 3}
	s.Add(1)
	fmt.Println(s)

	type User struct{ Name string }
	m := map[int]User{1: {"Tom"}}
	fmt.Println(m)
}
