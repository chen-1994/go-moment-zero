package slice_map_test

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

func TestMap1() {
	var m1 = map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
	fmt.Println(m1)

	m2 := make(map[string]int)
	m2["1"] = 1
	fmt.Println(m2)

	if val, exist := m1[3]; exist {
		fmt.Println(val)
	} else {
		fmt.Println("not key")
	}
	if _, exist := m1[2]; exist {
		m1[2] = "3"
	}
	fmt.Println(m1)

	//math.NaN() 特殊值
	mp := make(map[float64]string, 10)
	mp[math.NaN()] = "a"
	mp[math.NaN()] = "b"
	mp[math.NaN()] = "c"
	_, exist := mp[math.NaN()]
	fmt.Println(exist)
	fmt.Println(mp)

	//删除
	delete(m1, 2)
	fmt.Println(m1)
	///math.NaN() 无法删除
	delete(mp, math.NaN())
	fmt.Println(mp)

	//遍历
	for key, val := range m1 {
		fmt.Println(key, val)
	}
	//NaN 虽然无法正常取值，但可以遍历
	for key, val := range mp {
		fmt.Println(key, val)
	}

	//清空map delete 、 clear
	for key, _ := range m1 {
		delete(m1, key)
	}
	fmt.Println(m1)

	clear(mp)
	fmt.Println(mp)

	//set go并没有set 一般用map实现 使用struct因为值占用 0 字节
	set := make(map[int]struct{}, 10)
	for i := 0; i < 10; i++ {
		set[rand.Intn(100)] = struct{}{}
	}
	fmt.Println(set)
}

func MapTest2() {
	var group sync.WaitGroup

	// 1. 添加计数（假设要跑 10 个任务）
	group.Add(10)
	// map
	//mp := make(map[string]int, 10)
	var mp sync.Map //不能使用make sync.Map是结构体
	for i := 0; i < 10; i++ {
		go func() {
			// 写操作
			for i := 0; i < 100; i++ {
				//mp["helloworld"] = 1
				mp.Store("helloworld", 1)
			}
			// 读操作
			for i := 0; i < 10; i++ {
				//fmt.Println(mp["helloworld"])
				mp.Load("helloworld")
			}
			// 3. 每个任务结束时调用 Done (相当于 Add(-1))
			defer group.Done()
		}()
	}
	// 2. 阻塞在这里，直到计数器归零
	group.Wait()
}
