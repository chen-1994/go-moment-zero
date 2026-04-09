package slice_map_test

import (
	"fmt"
	"slices"
)

func Test() {
	// 数组
	var a [5]int
	a[0] = 10
	fmt.Println(a)
	a = [5]int{1, 2, 3, 4, 5}
	fmt.Println(a)
	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println(b)
	c := [5]int{1, 2, 3}
	fmt.Println(c)
	d := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(d)

	e := a[1:3]
	fmt.Println(e)
	f := a[:3]
	fmt.Println(f)
	g := a[1:]
	fmt.Println(g)
	h := a[:]
	fmt.Println(h)
}

func Test2() {
	a := []int{1, 2, 3, 4} //切片
	b := a[:]              //
	b[0] = 0
	fmt.Println(a, b)
	//没有clone 切片会改变原来值
	x := []int{1, 2, 3, 4}
	y := slices.Clone(x[:])
	y[0] = 0
	fmt.Println(x, y)
}

func Test3() {
	a := make([]int, 0, 0)
	a = append(a, 1, 2, 3, 4, 5, 6, 7)
	fmt.Println(len(a), cap(a))

	//头部插入
	a = append([]int{-1, 0}, a...)
	fmt.Println(a)
	//中间插入
	i := 2
	a = append(a[:i], append([]int{99, 100}, a[i:]...)...)
	fmt.Println(a)
	//尾部插入
	a = append(a, []int{9999}...)
	a = append(a, 1, 2, 3)
	fmt.Println(a)

	//头部删除
	a = a[1:]
	fmt.Println(a)
	//中间删除
	i = 2
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
	//尾部删除
	a = a[:len(a)-1]
	fmt.Println(a)

	//删除所有
	//a = a[:0]
	//fmt.Println(a)

	//copy
	b := make([]int, 0)
	fmt.Println(a, b)
	fmt.Println(copy(a, b))
	fmt.Println(a, b)
	c := make([]int, 10)
	fmt.Println(copy(c, a)) //不要写反了，否则反向copy
	fmt.Println(a, c)

	//遍历
	for i = 0; i < len(a); i++ {
		fmt.Println(a[i])
	}

	for idx, val := range a {
		fmt.Println(idx, val)
	}
}

func Test4() {
	//多维数组
	//数组的一维二维是固定的
	var a [5][5]int
	for _, val := range a {
		fmt.Println(val)
	}
	var b = make([][]int, 5)
	for _, val := range b {
		fmt.Println(val)
	}
	//切片的长度不是固定的，所以要单独初始化
	for i := 0; i < 5; i++ {
		b[i] = make([]int, 5)
	}
	for _, val := range b {
		fmt.Println(val)
	}
}

func Test5() {
	//拓展表达式
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // cap = 9
	s2 := s1[3:4]                          // cap = 9 - 3 = 6
	s2 = append(s2, 1)                     // 添加新元素，由于容量为6.所以没有扩容，直接修改底层数组
	fmt.Println(s2)
	fmt.Println(s1)

	s11 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // cap = 9
	s22 := s11[3:4:4]                       // cap = 4 - 3 = 1
	s22 = append(s22, 1)                    // 容量不足，分配新的底层数组
	fmt.Println(s22)
	fmt.Println(s11)

	//清理数据
	//clear 所有值设为0
	clear(s1)
	fmt.Println(s1)
	//清空切片
	s2 = s2[:0:0]
	fmt.Println(s2)
}
