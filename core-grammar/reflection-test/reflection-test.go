package reflection_test

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

//反射是一种在运行时检查语言自身结构的机制，它可以很灵活的去应对一些问题，但同时带来的弊端也很明显，例如性能问题等等。在 Go 中，
//反射与interface{}密切相关，很大程度上，只要有interface{}出现的地方，就会有反射。Go 中的反射 API 是由标准库reflect包提供的

// 接口
// 在 Go 中，接口本质上是结构体，Go 在运行时将接口分为了两大类，一类是没有方法集的接口，另一个类则是有方法集的接口。
// 对于含有方法集的接口来说，在运行时由如下的结构体iface来进行表示
//
//	type iface struct {
//	  tab  *itab // 包含 数据类型，接口类型，方法集等
//	  data unsafe.Pointer // 指向值的指针
//	}
//
// 而对于没有方法集接口来说，在运行时由eface 结构体来进行表示，如下
//
//	type eface struct {
//	  _type *_type // 类型
//	  data  unsafe.Pointer // 指向值的指针
//	}

// 而这两个结构体在reflect包下都有与其对应的结构体类型，iface对应的是nonEmptyInterface
//
//	type nonEmptyInterface struct {
//	 itab *struct {
//	   ityp *rtype // 静态接口类型
//	   typ  *rtype // 动态具体类型
//	   hash uint32 // 类型哈希
//	   _    [4]byte
//	   fun  [100000]unsafe.Pointer // 方法集
//	 }
//	 word unsafe.Pointer // 指向值的指针
//	}
//
// 而eface对应的是emptyInterface
//
//	type emptyInterface struct {
//	  typ  *rtype // 动态具体类型
//	  word unsafe.Pointer // 指向指针的值
//	}

// Go 语言是一个百分之百的静态类型语言，静态这一词是体现在对外表现的抽象的接口类型是不变的，而动态表示是接口底层存储的具体实现的类型是可以变化的

// 桥梁
// 在reflect包下，分别有reflect.Type接口类型来表示 Go 中的类型，reflect.Value结构体类型来表示 Go 中的值
// Go 中所有反射相关的操作都是基于这两个类型，reflect包提供了两个函数来将 Go 中的类型转换为上述的两种类型以便进行反射操作

// 核心
// 在 Go 中有三个经典的反射定律，结合上面所讲的内容也就非常好懂，分别如下
// 反射可以将interface{}类型变量转换成反射对象
// 反射可以将反射对象还原成interface{}类型变量
// 要修改反射对象，其值必须是可设置的
// 这三个定律便是 Go 反射的核心，当需要访问类型相关信息时，便需要用到reflect.TypeOf，当需要修改反射值时，就需要用到reflect.ValueOf 类型

func Test() {
	//类型
	str := "hello world"
	reflectType := reflect.TypeOf(str)
	fmt.Println(reflectType)

	var eface any
	eface = map[string]int{}
	rType := reflect.TypeOf(eface)
	fmt.Println(rType.Key().Kind())
	fmt.Println(rType.Elem().Kind())

	var eface2 any
	eface2 = new(strings.Builder) //// 赋值指针
	rType2 := reflect.TypeOf(eface2)
	VType2 := rType2.Elem()
	fmt.Println(VType2.PkgPath())
	fmt.Println(VType2.Name())
}

func Test1() {
	// size
	//通过Size方法可以获取对应类型所占的字节大小
	//使用unsafe.Sizeof()可以达到同样的效果
	fmt.Println(reflect.TypeOf(0).Size())
	fmt.Println(reflect.TypeOf("").Size())
	fmt.Println(reflect.TypeOf(complex(0, 0)).Size())
	fmt.Println(reflect.TypeOf(0.1).Size())
	fmt.Println(reflect.TypeOf([]string{}).Size())
}

// Comparable
// 通过Comparable方法可以判断一个类型是否可以被比较 意味着该类型的两个变量可以放在一起比对是否相等
func Test2() {
	fmt.Println(reflect.TypeOf("hello world!").Comparable())
	fmt.Println(reflect.TypeOf(1024).Comparable())
	fmt.Println(reflect.TypeOf([]int{}).Comparable())
	fmt.Println(reflect.TypeOf(struct{}{}).Comparable())
}

// Implements
func Test3() {
	//通过Implements方法可以判断一个类型是否实现了某一接口
	rIface := reflect.TypeOf(new(MyInterface)).Elem()
	fmt.Println(reflect.TypeOf(new(MyStruct)).Elem().Implements(rIface))
	fmt.Println(reflect.TypeOf(new(HisStruct)).Elem().Implements(rIface))
}

type MyInterface interface {
	My() string
}
type MyStruct struct {
}

func (MyStruct) My() string {
	return "hello world"
}

type HisStruct struct{}

func (h HisStruct) String() string {
	return "his"
}

// ConvertibleTo
// 使用ConvertibleTo方法可以判断一个类型是否可以被转换为另一个指定的类型
func Test4() {
	rIface := reflect.TypeOf(new(MyInterface)).Elem()
	fmt.Println(reflect.TypeOf(new(MyStruct)).Elem().ConvertibleTo(rIface))
	fmt.Println(reflect.TypeOf(new(HisStruct)).Elem().ConvertibleTo(rIface))
}

// ====值
func Test5() {
	str := "hello world"
	reflectValue := reflect.ValueOf(str)
	fmt.Println(reflectValue)

	// Type方法可以获取一个反射值的类型
	fmt.Println(reflectValue.Type())

	// Elem 获取一个反射值的元素反射值
	num := new(int)
	*num = 112233 //// 以指针为例子
	rValue := reflect.ValueOf(num).Elem()
	fmt.Println(rValue.Interface())
}

// 指针
//获取一个反射值的指针方式有两种
// 返回一个表示v地址的指针反射值
//func (v Value) Addr() Value

// 返回一个指向v的原始值的uinptr 等价于 uintptr(Value.Addr().UnsafePointer())
//func (v Value) UnsafeAddr() uintptr

// 返回一个指向v的原始值的uintptr
// 仅当v的Kind为 Chan, Func, Map, Pointer, Slice, UnsafePointer时，否则会panic
//func (v Value) Pointer() uintptr

// 返回一个指向v的原始值的unsafe.Pointer
// 仅当v的Kind为 Chan, Func, Map, Pointer, Slice, UnsafePointer时，否则会panic
// func (v Value) UnsafePointer() unsafe.Pointer
func Test6() {
	num := new(int)
	ele := reflect.ValueOf(&num).Elem()
	fmt.Println("&num", num)
	fmt.Println("Addr", ele.Addr())
	fmt.Println("UnsafeAddr", unsafe.Pointer(ele.UnsafeAddr()))
	fmt.Println("Pointer", unsafe.Pointer(ele.Addr().Pointer()))
	fmt.Println("UnsafePointer", ele.Addr().UnsafePointer())
	//fmt.Println会反射获取参数的类型，如果是reflect.Value类型的话，会自动调用Value.Interface()来获取其原始值。

	//换成一个 map 再来一遍
	dic := map[string]int{}
	elem := reflect.ValueOf(&dic).Elem()
	println(dic)
	fmt.Println("Addr", elem.Addr())
	fmt.Println("UnsafeAddr", *(*unsafe.Pointer)(unsafe.Pointer(elem.UnsafeAddr())))
	fmt.Println("Pointer", unsafe.Pointer(elem.Pointer()))
	fmt.Println("UnsafePointer", elem.UnsafePointer())
}

func Test7() {
	//// 设置值 func (v Value) Set(x Value)
	num := new(int)
	*num = 112233
	rValue := reflect.ValueOf(num)
	ele := rValue.Elem()
	fmt.Println(ele.Interface())
	ele.SetInt(100)
	fmt.Println(ele.Interface())
	//获取值 func (v Value) Interface() (i any)
	if v, ok := rValue.Interface().(*int); ok {
		fmt.Println(*v)
	}
}

// ===============函数
// 通过反射可以获取函数的一切信息，也可以反射调用函数
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func Test8() {
	// 信息
	rType := reflect.TypeOf(Max)
	fmt.Println(rType.Name())                  //// 输出函数名称,字面量函数的类型没有名称
	fmt.Println(rType.NumIn(), rType.NumOut()) //// 输出参数，返回值的数量
	rParamType := rType.In(0)
	fmt.Println(rParamType.Kind()) //// 输出第一个参数的类型
	rResType := rType.Out(0)
	fmt.Println(rResType.Kind()) // 输出第一个返回值的类型

	//调用
	// 获取函数的反射值
	rVal := reflect.ValueOf(Max)
	// 传入参数数组
	rResVal := rVal.Call([]reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)})
	for _, value := range rResVal {
		fmt.Println(value.Interface())
	}
}

// ======结构体
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	money   int
}

func (p Person) Talk(msg string) string {
	return msg
}

// 访问字段
//reflect.StructField结构的结构如下
//type StructField struct {
//  // 字段名称
//  Name string
//  // 包名
//  PkgPath string
//  // 类型名
//  Type      Type
//  // Tag
//  Tag       StructTag
//  // 字段的字节偏移
//  Offset    uintptr
//  // 索引
//  Index     []int
//  // 是否为嵌套字段
//  Anonymous bool
//}

// 访问结构体字段的方法有两种，一种是通过索引来进行访问，另一种是通过名称。
func Test9() {
	// 索引
	rType := reflect.TypeOf(new(Person)).Elem()
	fmt.Println(rType.NumField()) //// 输出结构体字段的数量
	for i := 0; i < rType.NumField(); i++ {
		structField := rType.Field(i)
		fmt.Println(structField.Index, structField.Name, structField.Type, structField.Offset, structField.IsExported())
	}
	// 名称
	if field, ok := rType.FieldByName("Name"); ok {
		fmt.Println(field.Name, field.Type, field.IsExported())
	}
}

// 修改字段
func Test10() {
	rValue := reflect.ValueOf(&Person{
		Name:    "",
		Age:     0,
		Address: "",
		money:   0,
	}).Elem()
	//获取字段
	name := rValue.FieldByName("Name")
	//修改字段值
	if (name != reflect.Value{}) {
		name.SetString("hello world")
	}
	fmt.Println(rValue.Interface())
	money := rValue.FieldByName("money")
	if (money != reflect.Value{}) {
		//构造指向该结构体为导出字段的指针反设置
		p := reflect.NewAt(money.Type(), money.Addr().UnsafePointer())
		field := p.Elem()
		field.SetInt(100)
	}
	fmt.Printf("%+v\n", rValue.Interface())
}

// 访问tag
// 获取到StructField后，便可以直接访问其 Tag
// // 如果不存在，ok为false
// func (tag StructTag) Lookup(key string) (value string, ok bool)
//
// // 如果不存在，返回空字符串
// func (tag StructTag) Get(key string) string
func Test11() {
	rType := reflect.TypeOf(new(Person)).Elem()
	name, ok := rType.FieldByName("Name")
	if ok {
		fmt.Println(name.Tag.Lookup("json"))
		fmt.Println(name.Tag.Get("json"))
	}
}

// 访问方法
// 访问方法与访问字段的过程很相似，只是函数签名略有区别。reflect.Method结构体如下
//
//	type Method struct {
//	 // 方法名
//	 Name string
//	 // 包名
//	 PkgPath string
//	 // 方法类型
//	 Type  Type
//	 // 方法对应的函数，第一个参数是接收者
//	 Func  Value
//	 // 索引
//	 Index int
//	}
func Test12() {
	rType := reflect.TypeOf(new(Person)).Elem()
	fmt.Println(rType.NumMethod()) //输出方法数
	for i := 0; i < rType.NumMethod(); i++ {
		method := rType.Method(i)
		fmt.Println(method.Index, method.Name, method.Type, method.IsExported())
		fmt.Println("方法参数")
		for j := 0; j < method.Type.NumIn(); j++ {
			fmt.Println(method.Type.In(j).String())
		}
		fmt.Println("方法返回值")
		for j := 0; j < method.Type.NumOut(); j++ {
			fmt.Println(method.Type.Out(j).String())
		}
	}

}

// 调用方法
// 调用方法与调用函数的过程相似，而且并不需要手动传入接收者
func Test13() {
	rValue := reflect.ValueOf(new(Person)).Elem()
	fmt.Println(rValue.NumField())
	talk := rValue.FieldByName("Talk")
	if (talk != reflect.Value{}) {
		result := talk.Call([]reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)})
		for _, value := range result {
			fmt.Println(value.Interface())
		}
	}
}

// ===============创建
// 通过反射可以构造新的值，reflect包同时根据一些特殊的类型提供了不同的更为方便的函数。
// 基本类型
func Test14() {
	rValue := reflect.New(reflect.TypeOf(*new(string)))
	rValue.Elem().SetString("hello world")
	fmt.Println(rValue.Elem().Interface())
}

// 结构体
func Test15() {
	rValue := reflect.TypeOf(new(Person)).Elem()
	person := reflect.New(rValue).Elem()
	fmt.Println(person.Interface())
}

// 切片
func Test16() {
	rValue := reflect.MakeSlice(reflect.TypeOf(*new([]int)), 10, 10)
	for i := 0; i < 10; i++ {
		rValue.Index(i).SetInt(int64(i))
	}
	fmt.Println(rValue.Interface())
}

// map
func Test17() {
	rValue := reflect.MakeMap(reflect.TypeOf(*new(map[string]int)))
	//rValue := reflect.MakeMapWithSize(reflect.TypeOf(*new(map[string]int)), 10)
	rValue.SetMapIndex(reflect.ValueOf("name"), reflect.ValueOf(2))
	fmt.Println(rValue.Interface())
}

// 管道
func Test18() {
	makeChan := reflect.MakeChan(reflect.TypeOf(new(chan int)).Elem(), 0)
	fmt.Println(makeChan.Interface())
}

// 函数
func Test19() {
	fn := reflect.MakeFunc(reflect.TypeOf(new(func(int))).Elem(), func(args []reflect.Value) (results []reflect.Value) {
		for _, arg := range args {
			fmt.Println(arg.Interface())
		}
		return nil
	})
	fmt.Println(fn.Type())
	fn.Call([]reflect.Value{reflect.ValueOf(1024)})
}

// 完全相等
// reflect.DeepEqual是反射包下提供的一个用于判断两个变量是否完全相等的函数，签名如下。
// func DeepEqual(x, y any) bool
// 该函数对于每一种基础类型都做了处理，下面是一些类型判断方式。
//
// 数组：数组中的每一个元素都完全相等
// 切片：都为nil时，判为完全相等，或者都不为空时，长度范围内的元素完全相等
// 结构体：所有字段都完全相等
// 映射表：都为nil时，为完全相等，都不为nil时，每一个键所映射的值都完全相等
// 指针：指向同一个元素或指向的元素完全相等
// 接口：接口的具体类型完全相等时
// 函数：只有两者都为nil时才是完全相等，否则就不是完全相等
func Test20() {
	//切片
	a := make([]int, 10)
	b := make([]int, 10)
	fmt.Println(reflect.DeepEqual(a, b))

	// 结构体
	mike := Person{
		Name:    "mike",
		Age:     20,
		Address: "mike",
	}
	jack := Person{
		Name:    "jack",
		Age:     20,
		Address: "jack",
	}
	tom := Person{
		Name:    "tom",
		Age:     20,
		Address: "tom",
	}
	fmt.Println(reflect.DeepEqual(mike, jack))
	fmt.Println(reflect.DeepEqual(tom, jack))
	fmt.Println(reflect.DeepEqual(jack, jack))
}
