package iter_test

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"slices"
	"testing"

	"github.com/saylorsolutions/x/iterx"
)

// 迭代器：用于迭代特定数据结构的关键字为for range
//仅能作用于语言内置的几个数据结构：数组 切片 字符串 map chan 整型值
//go1.23 版本更新以后，for range关键字支持了range over func，这样一来自定义迭代器也就成为了可能

// 闭包的斐波那契数列
func Fibonacci(n int) func() (int, bool) {
	a, b, c := 1, 1, 2
	i := 0
	return func() (int, bool) {
		if i > n {
			return 0, false
		} else if i < 2 {
			f := i
			i++
			return f, true
		}
		a, b = b, c
		c = a + b
		i++
		return a, true
	}
}

// 推送式迭代器
func Fibonacci2(n int) func(yield func(int) bool) {
	a, b, c := 0, 1, 1
	return func(yield func(int) bool) {
		for range n {
			if !yield(a) {
				return
			}
			a, b = b, c
			c = a + b
		}
	}
}
func myIter2(yield func(int, string) bool) {
	yield(0, "A")
	yield(1, "B")
}

func Test() {
	n := 8
	for f := range Fibonacci2(n) {
		fmt.Println(f)
	}
	//根据官方定义，上面迭代器Backward的例子使用就等同于下面这段代码
	//循环体的 body 就是迭代器的回调函数yiled，当函数返回true迭代器会继续迭代，否则就会停止。
	Fibonacci2(n)(func(f int) bool {
		fmt.Println(f)
		return true
	})

	//type Seq[V any] func(yield func(V) bool)
	// 1. 定义迭代器
	myIter := func(yield func(string) bool) {
		data := []string{"Apple", "Banana", "Cherry"}
		for _, v := range data {
			// 调用 yield 发送值给 for 循环
			// 如果 for 循环里执行了 break，yield 会返回 false
			if !yield(v) {
				return
			}
		}
	}
	// 2. 补全循环体
	for v := range myIter {
		fmt.Println("当前值:", v)
		if v == "Banana" {
			break // 此时 yield 会返回 false，上面的循环会停止
		}
	}
	//type Seq2[K, V any] func(yield func(K, V) bool)
	for i, v := range myIter2 {
		fmt.Printf("索引: %d, 值: %s\n", i, v)
	}

	//标准库中没有定义 0 个参数的 Seq，但这也是完全允许的，它相当于 func(yield func() bool)
}

// iterator 补全实现
func iterator() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s"}

		for i, v := range data {
			// 调用 yield 将 index 和 value 传给 for 循环
			// 如果 for 循环中途执行了 break，yield 会返回 false
			if !yield(i, v) {
				return // 终止迭代
			}
		}
	}
}

// 既然循环体中的代码是作为回调函数传入迭代器的，而且它很可能是一个闭包函数，
// Go 就需要让一个闭包函数在执行defer，return，break，goto等关键字时表现的像一个普通循环体代码段一样，思考下面几种情况。
func Test2() {
	//for index, value := range iterator() {
	//	fmt.Println(index, value)
	//}
	//比如说在迭代器循环中返回，那么在yield回调函数中要怎么去处理这个 return 呢？
	//不可能直接在回调函数中 return，这么做只会让迭代停止而已，达不到返回的效果
	for index, value := range iterator() {
		if index > 10 {
			return
		}
		fmt.Println(index, value)
	}
	//再比如说在迭代器循环中使用defer
	//也不能直接在回调函数中使用defer，因为这么做的话在回调函数结束时就会直接延迟调用了
	for index, value := range iterator() {
		defer fmt.Println(index, value)
	}
	//像其他的几个关键字break，continue，goto也是类似的，好在这些情况 Go 已经帮我们处理好了，我们只需使用即可，可以暂时不需要关心这些，如果感兴趣可以自行浏览
	//https://github.com/golang/go/blob/go1.23.0/src/cmd/compile/internal/rangefunc/rewrite.go#L628中的源代码。
}

// ======拉取式迭代器
// 推送式迭代器（pushing iterator）是由迭代器来控制迭代的逻辑，用户被动获取元素，
// 相反的拉取式迭代器（pulling iterator）就是由用户来控制迭代逻辑，主动的去获取序列元素。一般而言，
// 拉取式迭代器都会有特定的函数如next()，stop()来控制迭代的开始或结束，它可以是一个闭包或者结构体
func Test3() {
	file, err := os.Open("d://ads.properties")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 注意：scanner.Text() 已经处理了读取，
		// scanner.Err() 通常在循环结束后统一检查，但在循环内检查也是安全的
		line := scanner.Text() //Scanner 通过方法Text()来获取文件中的下一行文本,通过方法Scan()来表示迭代是否结束，这也是拉取式迭代器的一种模式
		if err := scanner.Err(); err != nil {
			fmt.Println("读取中出错", err)
		}
		fmt.Println(line)
	}

	//Scanner 采用结构体来记录状态，而在iter库定义的拉取式迭代器采用闭包来记录状态，
	//我们通过iter.Pull或iter.Pull2函数就可以将一个标准的推送式迭代器转换为拉取式迭代器，iter.Pull与iter.Pull2的区别就是后者的返回值有两个
	n := 10
	next, stop := iter.Pull(Fibonacci2(n))
	defer stop()
	for {
		fibn, ok := next()
		if !ok {
			break
		}
		fmt.Println(fibn)
	}
	// 闭包实现
	fib := Fibonacci(n)
	for {
		n, ok := fib()
		if !ok {
			fmt.Println(n)
		}
	}
	//转换过程：闭包 → 迭代器 → 拉取式迭代器，闭包与拉取式迭代器的用法都大差不差，它们的思想都是一样的，后者还会因为各种各样的处理导致性能上的拖累。
	//老实说这么做确实多此一举，它的应用场景确实不是很多，不过iter.pull是为了iter.Seq而存在的，
	//也就是为了将推送式迭代器转换成拉取式迭代器的而存在的，如果你仅仅只是想要一个拉取式迭代器，还专门为此去实现一个推送式迭代器来进行转换，
	//要这样做的话不妨考虑下自己实现的复杂度和性能，就像这个斐波那契数列的例子一样，绕了一圈又回到原点，唯一的好处可能就是符合官方的迭代器规范
}

// =======错误处理
// 在迭代时发生了错误怎么办？我们可以将其传递给yield函数让for range返回，让调用者来进行处理
// 值得注意的是，ScanLines迭代器是一次性使用的，文件关闭以后就不能再次使用了
func ScanLines(reader io.Reader) iter.Seq2[string, error] {
	scanner := bufio.NewScanner(reader)
	return func(yield func(string, error) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text(), scanner.Err()) {
				return
			}
		}
	}
}

func Test4() {
	file, err := os.Open("d://ads.properties")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果发生了 panic，就像平常一样使用recovery即可。
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", err)
			os.Exit(1)
		}
	}()
	for line, err := range ScanLines(file) {
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(line)
	}
	//这样处理起来就跟普通的错误处理没什么区别，拉取式迭代器也是同理
	next, stop := iter.Pull2(ScanLines(file))
	defer stop()
	for {
		line, err, ok := next()
		if err != nil {
			fmt.Println(err)
			break
		} else if !ok {
			break
		}
		fmt.Println(line)
	}
	//如果发生了 panic，就像平常一样使用recovery即可。拉取式迭代器依然同理，这里就不演示了。
}

// =========标准库
// 有很多标准库也支持了迭代器，最常用的就是slices和maps标准库
func Test5() {
	s := []int{1, 2, 3, 4, 5}
	//slices.All
	for i, n := range slices.All(s) {
		fmt.Println(i, n)
	}
	//slices.Values
	for n := range slices.Values(s) {
		fmt.Println(n)
	}
	//slices.Chunk函数会返回一个迭代器，该迭代器会以 n 个元素为切片推送给调用者 如下：n=2
	for chuck := range slices.Chunk(s, 2) {
		fmt.Println(chuck)
	}
	//slices.Collect函数会将切片迭代器收集成一个切片
	s2 := slices.Collect(slices.Values(s))
	fmt.Println(s2)

	m := map[string]int{"1": 1, "2": 2, "3": 3}
	//maps.Keys会返回一个迭代 map 所有键的迭代器，配合slices.Collect可以直接收集成一个切片
	keys := slices.Collect(maps.Keys(m))
	fmt.Println(keys)
	vals := slices.Collect(maps.Values(m))
	fmt.Println(vals)

	//maps.All可以将一个 map 转换为成一个 map 迭代器
	for k, v := range maps.All(m) {
		fmt.Println(k, v)
	}
	//maps.Collect可以将一个 map 迭代器收集成一个 map
	m2 := maps.Collect(maps.All(m))
	fmt.Println(m2)
	//collect 函数一般作为数据流处理的终结函数来使用
}

// ==================链式调用
// 通过上面标准库提供的函数，我们可以将其组合来处理数据流，比如对数据流进行排序，如下
// sortedSlices := slices.Sorted(slices.Values(s))
// go 的迭代器采用的是闭包，只能像这样嵌套函数调用，本身没法链式调用，调用链长了以后可读性会很差，但我们可以自己通过结构体来记录迭代器，就能够实现链式调用
type SliceSeq[E any] struct {
	seq iter.Seq2[int, E]
}

func (s SliceSeq[E]) All() iter.Seq2[int, E] {
	return s.seq
}
func (s SliceSeq[E]) Filter(filter func(int, E) bool) SliceSeq[E] {
	return SliceSeq[E]{
		seq: func(yield func(int, E) bool) {
			//重新组织索引
			i := 0
			for k, v := range s.seq {
				if filter(k, v) {
					if !yield(i, v) {
						return
					}
					i++
				}
			}
		},
	}
}

func (s SliceSeq[E]) Map(f func(int, E) E) SliceSeq[E] {
	return SliceSeq[E]{
		seq: func(yield func(int, E) bool) {
			for k, v := range s.seq {
				if !yield(k, v) {
					return
				}
			}
		},
	}
}

func (s SliceSeq[E]) Fill(fill E) SliceSeq[E] {
	return SliceSeq[E]{
		seq: func(yield func(int, E) bool) {
			for k, _ := range s.seq {
				if !yield(k, fill) {
					return
				}
			}
		},
	}
}
func (s SliceSeq[E]) Find(equla func(int, E) bool) (_ E) {
	for i, v := range s.seq {
		if equla(i, v) {
			return v
		}
	}
	return
}
func (s SliceSeq[E]) Some(match func(int, E) bool) bool {
	for i, v := range s.seq {
		if match(i, v) {
			return true
		}
	}
	return false
}
func (s SliceSeq[E]) Every(match func(int, E) bool) bool {
	for i, v := range s.seq {
		if match(i, v) {
			return false
		}
	}
	return true
}
func (s SliceSeq[E]) Collect() []E {
	var res []E
	for _, v := range s.seq {
		res = append(res, v)
	}
	return res
}
func (s SliceSeq[E]) Sort(cmp func(x, y E) int) []E {
	collect := s.Collect()
	slices.SortFunc(collect, cmp)
	return collect
}
func (s SliceSeq[E]) SortStable(cmp func(x, y E) int) []E {
	collect := s.Collect()
	slices.SortStableFunc(collect, cmp)
	return collect
}
func Slice[S ~[]E, E any](s S) SliceSeq[E] {
	return SliceSeq[E]{seq: slices.All(s)}
}

func Test6() {
	s := []string{"a", "b", "c", "d"}
	// 处理元素值
	all := iterx.SliceMap(s)
	for i, v := range all {
		fmt.Println(i, v)
	}
	//寻找某一个指定值

}

var s []int

const n = 1000

func init() {
	for i := range n {
		s = append(s, i)
	}
}

func testNaiveFor(s []int) {
	for i, n := range s {
		_ = i
		_ = n
	}
}
func testPushing(s []int) {
	for i, n := range slices.All(s) {
		_ = i
		_ = n
	}
}
func testPulling(s []int) {
	next, stop := iter.Pull2(slices.All(s))
	for {
		i, m, ok := next()
		if !ok {
			stop()
			return
		}
		_ = i
		_ = m
	}
}
func BenchmarkNaive_10000(b *testing.B) {
	for range b.N {
		testNaiveFor(s)
	}
}
func BenchmarkPushing_10000(b *testing.B) {
	for range b.N {
		testPushing(s)
	}
}
func BenchmarkPulling_10000(b *testing.B) {
	for range b.N {
		testPulling(s)
	}
}

// 性能
func Test7() {
	b := &testing.B{N: n}
	BenchmarkNaive_10000(b)
	BenchmarkPushing_10000(b)
	BenchmarkPulling_10000(b)
}
