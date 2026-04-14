package generic_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

// 泛型
// 类型形参：T 就是一个类型形参，形参具体是什么类型取决于传进来什么类型
// 类型约束：int | float64 构成了一个类型约束，这个类型约束内规定了哪些类型是允许的，约束了类型形参的类型范围
func sum[T int | float64](a T, b T) T {
	return a + b
}

// 泛型切片，类型约束为 int | int32 | int64
type GenericSlice[T int | int32 | int64] []T

// 泛型哈希表，键的类型必须是可比较的，所以使用 comparable 接口，值的类型约束为 V int | string | byte
type GenericMap[K comparable, V int | string | byte] map[K]V

// 泛型结构体，类型约束为 T int | string
type GenericStruct[T int | string] struct {
	Name string
	id   T
}

// 泛型切片形参
type Company[T int | string, S []T] struct {
	Name  string
	Id    T
	Stuff S
}

// 泛型接口
type SayAble[T int | string] interface {
	Say() T
}

// fmt.Stringer 为接口 作为形参
func PrintObj[T fmt.Stringer](s T) {
	fmt.Println(s)
}

type Person struct {
	Name string
}

func (p Person) String() string {
	return fmt.Sprintf("%s", p.Name)
}

// 类型参数
func Write[W io.Writer](w W, bs []byte) (int, error) {
	return w.Write(bs)
}

// ============泛型断言
func Assert[T any](v any) (bool, T) {
	var av T
	if v == nil {
		return false, av
	}
	av, ok := v.(T)
	return ok, av
}

// 类型实参：Sum[int](1,2)，手动指定了 int 类型，int 就是类型实参
func Test() {
	//显式知名使用类型
	sum[int](1, 2)
	//不指定类型由编译器自行推断
	sum(1.1, 2.2)
	a := GenericSlice[int]{1, 2, 3}
	fmt.Println(a)

	gmap1 := GenericMap[int, string]{1: "hello world"}
	gmap2 := make(GenericMap[string, byte], 0)
	fmt.Println(gmap1, gmap2)

	gs1 := GenericStruct[string]{
		Name: "hello world",
		id:   "1",
	}
	gs2 := GenericStruct[int]{
		Name: "hello world",
		id:   2,
	}
	fmt.Println(gs1, gs2)

	c := Company[int, []int]{
		Name:  "hello world",
		Id:    1,
		Stuff: []int{1, 2, 3},
	}
	fmt.Println(c)

	PrintObj(Person{Name: "Alice"})

	str := "abc"
	n, err := Write(os.Stdout, []byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)

	_, v := Assert[string]("hello world")
	fmt.Println(v)
}

// ============类型集
// 在 1.18 以后，接口的定义变为了类型集 (type set)，含有类型集的接口又称为 General interfaces 即通用接口
// 类型集只能用于泛型中的类型约束，不能用作类型声明，类型转换，类型断言

type SignedInt interface {
	int8 | int16 | int | int32 | int64
}

type UnSignedInt interface {
	uint8 | uint16 | uint32 | uint64
}

// 并集 ： 取所有的
type Integer interface {
	SignedInt | UnSignedInt
}

type Integer2 interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

// 交集 ：取相交的
type Number interface {
	SignedInt
	Integer2
}

func DO[T Number](v T) T {
	return v
}

// 空集：没有交集
type Integer3 interface {
	SignedInt
	UnSignedInt
}

func Do1[T Integer3](v T) T {
	return v
}

// 空接口
// 空接口与空集并不同，空接口是所有类型集的集合，即包含所有类型
func Do2[T interface{}](v T) T {
	return v
}

// 底层类型
// 当使用 type 关键字声明了一个新的类型时，即便其底层类型包含在类型集内，当传入时也依旧会无法通过编译
type Int interface {
	int8 | int16 | int32 | int64
}
type TinyInt int8

func Do3[T Int](v T) T {
	return v
}

// 两种解决办法：
// 1:往类型集中并入该类型 但是这毫无意义，因为 TinyInt 与 int8 底层类型就是一致的
type Int1 interface {
	int8 | int16 | int | int32 | int64 | uint8 | uint16 | uint | uint32 | uint64 | TinyInt
}

// 2:使用 ~ 符号，来表示底层类型，如果一个类型的底层类型属于该类型集，那么该类型就属于该类型集
type Int2 interface {
	~int8 | ~int16 | ~int | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint | ~uint32 | ~uint64
}

func Test2() {

	DO[int8](1)
	//DO[uint8](1)//无法通过编译 因为交集没有uint8

	//Do1[int16](1) //空集，无法通过编译，传什么都不行
	Do2[int16](1)
	Do2[struct{}](struct{}{})

	//Do3[TinyInt](1) //无法通过编译，即便其底层类型属于Int类型集的范围内
}

//注意：泛型不能作为一个类型的基本类型
// 错误写法：type GenericType[T int | int32 | int64] T
//type GenericType[T int | int32 | int64] int //虽然下列的写法是允许的，不过毫无意义而且可能会造成数值溢出的问题，所以并不推荐

//泛型类型无法使用类型断言
//对泛型类型使用类型断言将会无法通过编译，泛型要解决的问题是 类型无关 的，如果一个问题需要根据不同类型做出不同的逻辑，那么就根本不应该使用泛型，应该使用 interface{} 或者 any。
//func Sum[T int | float64](a, b T) T {
//	ints,ok := a.(int) // 不被允许
//	switch a.(type) { // 不被允许
//	case int:
//	case bool:
//		...
//	}
//	return a + b
//}

//匿名结构不支持泛型 如下的代码将无法通过编译
//testStruct := struct[T int | string] {
//	Name string
//	Id T
//	}[int]{
//		Name: "jack",
//		Id: 1
//	}

//匿名函数不支持自定义泛型 以下两种写法都将无法通过编译
//var sum[T int | string] func (a, b T) T
//sum := func[T int | string](a,b T) T{
//	...
//}
//但是可以 使用 已有的泛型类型，例如闭包中
//func Sum[T int | float64](a, b T) T {
//  sub := func(c, d T) T {
//    return c - d
//  }
//  return sub(a,b) + a + b
//}

//不支持泛型方法
//类型集无法作为类型实参
//类型集中的交集问题
//类型集无法直接或间接的并入自身
//接口无法并入类型集
//方法集无法并入类型集

// =====================使用
// 队列
type Queue[T any] []T

func (q *Queue[T]) Push(e T) {
	*q = append(*q, e)
}
func (q *Queue[T]) Pop(e T) (_ T) {
	if q.Size() > 0 {
		res := q.peek()
		*q = (*q)[1:]
		return res
	}
	return
}
func (q *Queue[T]) peek() (_ T) { //_ T 作用是 return 等同于 "return T 类型的零值"
	if q.Size() > 0 {
		return (*q)[0]
	}
	return
}
func (q *Queue[T]) Size() int {
	return len(*q)
}

// 堆 可以在 O(1) 的时间内判断最大或最小值，所以它对元素有一个要求，
// 那就是必须是可以排序的类型，但内置的可排序类型只有数字和字符串，所以在堆的初始化时，需要传入一个自定义的比较器，比较器由调用者提供，并且比较器也必须使用泛型
type Comparator[T any] func(a, b T) int
type BinaryHeap[T any] struct {
	s []T
	c Comparator[T]
}

func (b *BinaryHeap[T]) Peek() (_ T) {
	if b.Size() > 0 {
		return b.s[0]
	}
	return
}

func (b *BinaryHeap[T]) Pop() (_ T) {
	size := b.Size()
	if size > 0 {
		res := b.s[0]
		b.s[0], b.s[size-1] = b.s[size-1], b.s[0]
		b.s = b.s[:size-1]
		b.down(0)
		return res
	}
	return
}

func (b *BinaryHeap[T]) Push(e T) {
	b.s = append(b.s, e)
	b.up(b.Size() - 1)
}

func (b *BinaryHeap[T]) up(i int) {
	if b.Size() == 0 || i <= 0 || i >= b.Size() {
		return
	}
	//for parentIndex := (i - 1) >> 1; parentIndex >= 0; parentIndex = parentIndex>>1 - 1 {
	//	// greater than or equal to
	//	if b.c(b.s[i], b.s[parentIndex]) >= 0 {
	//		break
	//	}
	//	b.s[i], b.s[parentIndex] = b.s[parentIndex], b.s[i]
	//	i = parentIndex
	//}
	for i > 0 {
		p := (i - 1) >> 1
		if b.c(b.s[i], b.s[p]) >= 0 {
			break
		}
		b.s[i], b.s[p] = b.s[p], b.s[i]
		i = p
	}
}

func (b *BinaryHeap[T]) down(i int) {
	if b.Size() == 0 || i < 0 || i >= b.Size() {
		return
	}
	size := b.Size()
	for lsonIndex := i<<1 + 1; lsonIndex < size; lsonIndex = i<<1 + 1 {
		rsonIndex := lsonIndex + 1

		if rsonIndex < size && b.c(b.s[rsonIndex], b.s[lsonIndex]) < 0 {
			lsonIndex = rsonIndex
		}
		// less than or equal to
		if b.c(b.s[i], b.s[lsonIndex]) <= 0 {
			break
		}
		b.s[i], b.s[lsonIndex] = b.s[lsonIndex], b.s[i]
		i = lsonIndex
	}
}

func (b *BinaryHeap[T]) Size() int {
	return len(b.s)
}

func NewHeap[T any](n int, c Comparator[T]) BinaryHeap[T] {
	var b BinaryHeap[T]
	b.s = make([]T, 0, n)
	b.c = c
	return b
}

func Test3() {

	var c Comparator[string] = func(a, b string) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	b := NewHeap[string](1, c)
	b.Push("hello")
	b.Push("world")
	b.Push("hi")
	b.Peek()
	b.Pop()
	fmt.Println(b)
}

// 对象池
type Pool[T any] struct {
	pool sync.Pool
}

func (p *Pool[T]) put(v T) {
	p.pool.Put(v)
}
func (p *Pool[T]) get() T {
	var v T
	get := p.pool.Get()
	if get != nil {
		v = get.(T)
	}
	return v
}

func NewPool[T any](newFn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFn()
			},
		},
	}
}

func Test4() {
	bufferPool := NewPool(func() *bytes.Buffer {
		return bytes.NewBuffer(nil)
	})
	for range 100 {
		buffer := bufferPool.get()
		buffer.WriteString("hello")
		fmt.Println(buffer.String())
		buffer.Reset()
		bufferPool.put(buffer)
	}
}
