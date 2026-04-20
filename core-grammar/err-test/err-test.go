package err_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"

	errors2 "github.com/pkg/errors"
)

// 异常三种级别
// error：正常的流程出错，需要处理，直接忽略掉不处理程序也不会崩溃
// panic：很严重的问题，程序应该在处理完问题后立即退出
// fatal：非常致命的问题，程序应该立即退出

func Test() {
	sum, err := checksum("./err-test/err-test.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sum)
}
func checksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	var haxSum [64]byte
	sum := hash.Sum(haxSum[:0])
	hex.Encode(haxSum[:], sum)
	return string(haxSum[:]), nil
}

// ===============error
// error 属于是一种正常的流程错误，它的出现是可以被接受的，大多数情况下应该对其进行处理，
// 当然也可以忽略不管，error 的严重级别不足以停止整个程序的运行。error本身是一个预定义的接口，该接口下只有一个方法Error()，该方法的返回值是字符串，用于输出错误信息
func Test2() {
	err1 := errors.New("test error")
	err2 := fmt.Errorf("test2: %w", err1)
	fmt.Println(err2)

}

// var定义的变量
var (
	ErrInvalid    = fs.ErrInvalid
	ErrPermission = fs.ErrPermission
	ErrExist      = fs.ErrExist
	ErrNotExist   = fs.ErrNotExist
	ErrClosed     = fs.ErrClosed

	ErrNoDeadline = errNoDeadline()
)

func errNoDeadline() error {
	return errors.New("test error")
}

// 自定义错误
func New(text string) error {
	return errors.New(text)
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// 传递
// 调用者调用的函数返回了一个错误，但是调用者本身不负责处理错误，于是也将错误作为返回值返回，抛给上一层调用者，这个过程叫传递，
// 错误在传递的过程中可能会层层包装，当上层调用者想要判断错误的类型来做出不同的处理时，可能会无法判别错误的类别或者误判，而链式错误正是为了解决这种情况而出现的
type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}
func (e *wrapError) Unwrap() error {
	return e.err
}
func Test3() {
	//wrappError同样实现了error接口，也多了一个方法Unwrap，用于返回其内部对于原 error 的引用，
	//层层包装下就形成了一条错误链表，顺着链表上寻找，很容易就能找到原始错误。由于该结构体并不对外暴露，所以只能使用fmt.Errorf函数来进行创建
	err := errors.New("test error")
	wrapError := fmt.Errorf("错误。%w", err) //使用时，必须使用%w格式动词，且参数只能是一个有效的 error
	fmt.Println(wrapError)
}

// 处理
type TimeError struct {
	Msg  string
	Time time.Time
}

func (e *TimeError) Error() string {
	return e.Msg
}
func NewMyError(err string) error {
	return &TimeError{
		Msg:  err,
		Time: time.Now(),
	}
}
func wrap1() error { // 包裹原始错误
	return fmt.Errorf("wrapp error: %w", wrap2())
}
func wrap2() error { // 原始错误
	return NewMyError("original error")
}
func Test4() {
	var myErr *TimeError
	err := wrap1()
	//检查错误链中是否有*TimeError类型的错误
	if errors.As(err, &myErr) {
		fmt.Println("myErr:", myErr.Time)
	}
}

// github.com/pkg/errors 增强包
// 通过格式化输出，就可以看到堆栈信息了，默认情况下是不会输出堆栈的。这个包相当于是标准库errors包的加强版，同样都是官方写的
func Do() error {
	return errors2.New("test error")
}
func Test5() {
	if err := Do(); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// ===================panic
// nil的 map 写入值的例子，肯定会触发 panic
func Test6() {
	var dic map[string]int
	dic["a"] = 'a' //panic: assignment to entry in nil map
}

//当程序中存在多个协程时，只要任一协程发生panic，如果不将其捕获的话，整个程序都会崩溃

// 创建
// 显式的创建panic十分简单，使用内置函数panic即可，函数签名如下:func panic(v any)
func Test7() {
	initDataBase("", 0)
}
func initDataBase(host string, port int) {
	if len(host) == 0 || port == 0 {
		panic("非法的数据链接参数")
	}
}

// 善后
// 程序因为panic退出之前会做一些善后工作，例如执行defer语句
func Test8() {
	defer fmt.Println("A")
	defer fmt.Println("B")
	fmt.Println("C")
	panic("panic")
	defer fmt.Println("D") //不会执行
}
func dangerOp() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	panic("panic")
	defer fmt.Println("3")
}

func Test9() {
	defer fmt.Println("A")
	defer fmt.Println("B")
	fmt.Println("C")
	dangerOp()
	defer fmt.Println("D")
}

// defer中也可以嵌套panic
func Test10() {
	defer fmt.Println("A")
	defer func() {
		func() {
			panic("panicA")
			defer fmt.Println("E")
		}()
	}()
	fmt.Println("C")
	dangerOp()
	defer fmt.Println("D")
	//综上所述，当发生panic时，会立即退出所在函数，并且执行当前函数的善后工作，例如defer，然后层层上抛，上游函数同样的也进行善后工作，直到程序停止运行。
}

// 当子协程发生panic时，不会触发当前协程的善后工作，如果直到子协程退出都没有恢复panic，那么程序将会直接停止运行。
func Test11() {
	demo()
	//当子协程发生panic时，父协程早已完成的函数的执行，进入了善后工作，在执行最后一个defer时，碰巧遇到了子协程发生panic，所以程序就直接退出运行
}
func demo() {
	defer func() {
		time.Sleep(time.Millisecond * 20)
		fmt.Println("A")
	}()
	fmt.Println("C")
	go dangerOp2()
	defer fmt.Println("D")
}
func dangerOp2() {
	time.Sleep(time.Millisecond)
	defer fmt.Println(1)
	defer fmt.Println(2)
	panic("panicB")
	defer fmt.Println(3)
}

// 恢复
func Test12() {
	dangerOp3()
	fmt.Println("程序正常退出")
}
func dangerOp3() {
	//当发生panic时，使用内置函数recover()可以及时的处理并且保证程序继续运行，必须要在defer语句中运行
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("recover")
		}
	}()
	panic("panic")
	//调用者完全不知道dangerOp()函数内部发生了panic，程序执行剩下的逻辑后正常退出，所以输出

	//recover()的使用有许多隐含的陷阱
	//1.闭包函数可以看作调用了一个函数，panic是向上传递而不是向下，自然闭包函数也就无法恢复panic
	//defer func() {
	//	func() {
	//		if err := recover(); err != nil {
	//			fmt.Println(err)
	//			fmt.Println("panic恢复")
	//		}
	//	}()
	//}()
	//panic("panic")

	//2.panic()的参数是nil
	//这种情况panic确实会恢复，但是不会输出任何的错误信息
	//panic(nil)
}

// fatal
// fatal是一种极其严重的问题，当发生fatal时，程序需要立刻停止运行，不会执行任何善后工作，通常情况下是调用os包下的Exit函数退出程序
func dangerOp4(str string) {
	if len(str) == 0 {
		fmt.Println("fatal")
		os.Exit(1)
	}
	fmt.Println()
}

func Test13() {
	dangerOp4("")
	//fatal级别的问题一般很少会显式的去触发，大多数情况都是被动触发
}
