package concurrency_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ======协程
// 协程（coroutine）是一种轻量级的线程，或者说是用户态的线程，不受操作系统直接调度，由 Go 语言自身的调度器进行运行时调度，因此上下文切换开销非常小
// 具有返回值的内置函数不允许跟随在 go 关键字后面
func Test1() {
	go fmt.Println("hello world1")
	go hello()
	go func() {
		fmt.Println("hello world3")
	}()
	//以上三种开启协程的方式都是可以的，但是其实这个例子执行过后在大部分情况下什么都不会输出，协程是并发执行的，
	//系统创建协程需要时间，而在此之前，主协程早已运行结束，一旦主线程退出，其他子协程也就自然退出了
}
func hello() {
	fmt.Println("hello world2")
}

//并发控制方法有三种：
//channel：管道
//WaitGroup：信号量
//Context：上下文

// 三种方法有着不同的适用情况，WaitGroup 可以动态的控制一组指定数量的协程，Context 更适合子孙协程嵌套层级更深的情况，管道更适合协程间通信。对于较为传统的锁控制，Go 也对此提供了支持：
// Mutex：互斥锁
// RWMutex ：读写互斥锁
func Test2() {
	intch := make(chan int, 1) //缓冲区为1的管道 也可以不设置缓冲区大小
	defer close(intch)         //关闭通道
	//写入数据
	intch <- 1
	//intch <- 2 //缓冲区大小为1，同步写入异常
	fmt.Println(<-intch)
	intch <- 2
	if ints, ok := <-intch; ok {
		fmt.Println(ints)
	}

	//无缓冲管道不应该同步的使用，正确来说应该开启一个新的协程来发送数据
	ch := make(chan int)
	defer close(ch)
	go func() {
		ch <- 1
	}()
	n := <-ch
	fmt.Println(n)
}

func Test3() {
	//一个有缓冲管道用于协程间通信，两个无缓冲管道用于同步父子协程的执行顺序
	ch := make(chan int, 5)
	chW := make(chan struct{})
	chR := make(chan struct{})
	defer func() {
		close(chR)
		close(ch)
		close(chW)
	}()
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			fmt.Println("写入", i)
		}
		chW <- struct{}{}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Microsecond)
			fmt.Println("读取", <-ch)
		}
		chR <- struct{}{}
	}()
	fmt.Println("写入完毕", <-chW)
	fmt.Println("读取完毕", <-chR)
}

// 通过内置函数 len 可以访问管道缓冲区中数据的个数，通过 cap 可以访问管道缓冲区的大小。
func Test4() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(len(ch), cap(ch))
}

// 利用管道的阻塞条件，可以很轻易的写出一个主协程等待子协程执行完毕的例
func Test5() {
	ch := make(chan struct{})
	defer close(ch)
	go func() {
		fmt.Println(2)
		ch <- struct{}{}
	}()
	<-ch
	fmt.Println(1)
}

// 通过有缓冲管道还可以实现一个简单的互斥锁，看下面的例子
var count = 0
var lock = make(chan struct{}, 1)

func add() {
	//加锁
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行加法")
	count++
	//解锁
	<-lock
}
func sub() {
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行减法")
	count--
	<-lock
}

// 单项通道
// 生产者：限定只能往里写，不能读
func produce(out chan<- int) {
	for i := 0; i < 10; i++ {
		out <- i
	}
	close(out)
}

// 消费者：限定只能从里读，不能关，也不能写
func consume(in <-chan int) {
	for v := range in {
		fmt.Println("接收到:", v)
	}
}

func Test6() {
	ch := make(chan int) // 创建一个双向通道
	go produce(ch)       // 隐式转换为 chan<- int
	consume(ch)          // 隐式转换为 <-chan int

	//关键规则
	//自动转换：你可以把一个「双向通道」赋值给「单向通道」变量。
	//不可逆转：你不能把一个「单向通道」转回「双向通道」，也不能在 chan<- 上进行接收操作。
	//关闭限制：
	//可以关闭 chan<- int（发送方关闭通道是标准做法）。
	//不能关闭 <-chan int（接收方关闭通道会导致编译错误，因为关闭操作本质上是一种发送行为）

	//常见坑点
	//如果你尝试在代码里直接定义一个单向通道而不初始化，它是没用的：
	//ch := make(chan<- int) // 这定义了一个只能写的通道，但由于没人能读，写进去就会永远阻塞
	//所以单向通道的唯一正确用法是作为函数签名的一部分，用来做权限控制。
}

// for range 遍历读取
func Test7() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		defer close(ch)
	}()
	/*for v := range ch {
		fmt.Println(v)
	}*/
	for i := 0; i < 11; i++ {
		v, ok := <-ch
		fmt.Println(v, ok)
	}
	//由于管道已经关闭了，即便缓冲区为空，再读取数据也不会导致当前协程阻塞，可以看到在第六次遍历的时候读取的是零值，并且 ok 为 false
	//关于管道关闭的时机，应该尽量在向管道发送数据的那一方关闭管道，而不要在接收方关闭管道，因为大多数情况下接收方只知道接收数据，并不知道该在什么时候关闭管道。

}

// WaitGroup
// sync.WaitGroup 是 sync 包下提供的一个结构体，WaitGroup 即等待执行，使用它可以很轻易的实现等待一组协程的效果。该结构体只对外暴露三个方法
// Add 方法用于指明要等待的协程的数量	func (wg *WaitGroup) Add(delta int)
// Done 方法表示当前协程已经执行完毕	func (wg *WaitGroup) Done()
// Wait 方法等待子协程结束，否则就阻塞	func (wg *WaitGroup) Wait()
// 当计数变为负数，或者计数数量大于子协程数量时，将会引发 panic。
func Test8() {
	var mainWait sync.WaitGroup
	var wait sync.WaitGroup
	mainWait.Add(10)
	fmt.Println("start")
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go func() {
			fmt.Println(i)

			wait.Done()
			mainWait.Done()
		}()
		wait.Wait()
	}
	mainWait.Wait()
	fmt.Println("end")
}

//WaitGroup 通常适用于可动态调整协程数量的时候，例如事先知晓协程的数量，又或者在运行过程中需要动态调整。WaitGroup 的值不应该被复制，
//复制后的值也不应该继续使用，尤其是将其作为函数参数传递时，应该传递指针而不是值。倘若使用复制的值，计数完全无法作用到真正的 WaitGroup 上，
//这可能会导致主协程一直阻塞等待，程序将无法正常运行。例如下方的代码
//func main() {
//  var mainWait sync.WaitGroup
//  mainWait.Add(1)
//  hello(mainWait)
//  mainWait.Wait()
//  fmt.Println("end")
//}
//func hello(wait sync.WaitGroup) {
//  fmt.Println("hello")
//  wait.Done()
//}

// Context
// Context 译为上下文，是 Go 提供的一种并发控制的解决方案，相比于管道和 WaitGroup，
// 它可以更好的控制子孙协程以及层级更深的协程。Context 本身是一个接口，
// 只要实现了该接口都可以称之为上下文例如著名 Web 框架 Gin 中的 gin.Context。
// context 标准库也提供了几个实现，分别是	emptyCtx、cancelCtx、timerCtx、valueCtx
var waitGroup sync.WaitGroup

// emptyCtx
//emptyCtx 就是空的上下文，context 包下所有的实现都是不对外暴露的，但是提供了对应的函数来创建上下文
//emptyCtx 通常是用来当作最顶层的上下文，在创建其他三种上下文时作为父上下文传入

// valueCtx
// valueCtx 实现比较简单，其内部只包含一对键值对，和一个内嵌的 Context 类型的字段
// 其本身只实现了 Value 方法，逻辑也很简单，当前上下文找不到就去父上下文找
func Test9() {
	waitGroup.Add(1)
	go Do(context.WithValue(context.Background(), 1, 2))
	waitGroup.Wait()
}
func Do(ctx context.Context) {
	//新建定时器
	ticker := time.NewTicker(time.Second)
	defer waitGroup.Done()
	for {
		select {
		case <-ctx.Done(): //永远不会执行
		case <-ticker.C:
			fmt.Println("timeout")
			return
		default:
			fmt.Println(ctx.Value(1))
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// cancelCtx
// cancelCtx 以及 timerCtx 都实现了 canceler 接口
// cancel 方法不对外暴露，在创建上下文时通过闭包将其包装为返回值以供外界调用
// cancelCtx 译为可取消的上下文，创建时，如果父级实现了 canceler，就会将自身添加进父级的 children 中，否则就一直向上查找。
// 如果所有的父级都没有实现 canceler，就会启动一个协程等待父级取消，然后当父级结束时取消当前上下文。
// 当调用 cancelFunc 时，Done 通道将会关闭，该上下文的任何子级也会随之取消，最后会将自身从父级中删除
func Test10() {
	bkg := context.Background()
	cancelCtx, cancel := context.WithCancel(bkg)
	waitGroup.Add(1)
	go func(ctx context.Context) {
		defer waitGroup.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				return
			default:
				fmt.Println("等待取消中...")
			}
			time.Sleep(time.Second)
		}
	}(cancelCtx)
	time.Sleep(time.Second)
	cancel()
	waitGroup.Wait()
}
func Test11() {
	waitGroup.Add(3)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go HttpHandle(ctx)
	time.Sleep(time.Second)
	cancelFunc()
	waitGroup.Wait()
}
func HttpHandle(ctx context.Context) {
	cancelCtxAuth, cancelAuth := context.WithCancel(ctx)
	cancelCtxMail, cancelMail := context.WithCancel(ctx)
	defer cancelAuth()
	defer cancelMail()
	defer waitGroup.Done()

	go AuthService(cancelCtxAuth)
	go MailService(cancelCtxMail)

	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		default:
			fmt.Println("正在处理http请求")
		}
		time.Sleep(time.Millisecond * 200)
	}
}
func AuthService(ctx context.Context) {
	defer waitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("auth 父级取消", ctx.Err())
			return
		default:
			fmt.Println("auth...")
		}
		time.Sleep(time.Millisecond * 200)
	}
}
func MailService(ctx context.Context) {
	defer waitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("mail 父级取消", ctx.Err())
			return
		default:
			fmt.Println("mail...")
		}
		time.Sleep(time.Millisecond * 200)
	}
}

// timerCtx
// timerCtx 在 cancelCtx 的基础之上增加了超时机制，context 包下提供了两种创建的函数，分别是 WithDeadline 和 WithTimeout，两者功能类似，
// 前者是指定一个具体的超时时间，比如指定一个具体时间 2023/3/20 16:32:00，后者是指定一个超时的时间间隔，比如 5 分钟后。两个函数的签名如下
// func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
// func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
// timerCtx 会在时间到期后自动取消当前上下文，取消的流程除了要额外的关闭 timer 之外，基本与 cancelCtx 一致
func Test12() {
	deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	waitGroup.Add(1)
	go func(ctx context.Context) {
		defer waitGroup.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("上下文取消", ctx.Err())
				return
			default:
				fmt.Println("等待取消中...")
			}
			time.Sleep(time.Millisecond * 200)
		}
	}(deadline)
	waitGroup.Wait()
	//WithTimeout 其实与 WithDealine 非常相似，它的实现也只是稍微封装了一下并调用 WithDeadline，和上面例子中的 WithDeadline 用法一样
	//就跟内存分配后不回收会造成内存泄漏一样，上下文也是一种资源，如果创建了但从来不取消，一样会造成上下文泄露，所以最好避免此种情况的发生。
}

// select
// select 在 Linux 系统中，是一种 IO 多路复用的解决方案，类似的，在 Go 中，select 是一种管道多路复用的控制结构。
// 什么是多路复用，简单的用一句话概括：在某一时刻，同时监测多个元素是否可用，被监测的可以是网络请求，文件 IO 等。
// 在 Go 中的 select 监测的元素就是管道，且只能是管道。select 的语法与 switch 语句类似
func Test13() {
	chA := make(chan int)
	chB := make(chan int)
	chC := make(chan int)
	defer func() {
		close(chA)
		close(chB)
		close(chC)
	}()

	l := make(chan struct{})

	go Send(chA)
	go Send(chB)
	go Send(chC)

	go func() {
	Loop:
		for {
			select {
			case n, ok := <-chA:
				fmt.Println("A", n, ok)
			case n, ok := <-chB:
				fmt.Println("B", n, ok)
			case n, ok := <-chC:
				fmt.Println("C", n, ok)
			case <-time.After(time.Second): //设置超时时间
				fmt.Println("timeout")
				break Loop
			}
		}
		l <- struct{}{}
	}()
	<-l
}
func Send(ch chan<- int) {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Millisecond)
		ch <- i
	}
}

// 超时
// time.After 函数，其返回值是一个只读的管道，该函数配合 select 使用可以非常简单的实现超时机制
func Test14() {
	chA := make(chan int)
	defer close(chA)
	go func() {
		time.Sleep(time.Second * 2)
		chA <- 1
	}()
	select {
	case n := <-chA:
		fmt.Println(n)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

// 永久阻塞
// 当 select 语句中什么都没有时，就会永久阻塞
func Test15() {
	fmt.Println("start")
	select {}
	fmt.Println("end")
}

// 在 select 的 case 中对值为 nil 的管道进行操作的话，并不会导致阻塞，该 case 则会被忽略，永远也不会被执行。例如下方代码无论执行多少次都只会输出 timeout。
func Test16() {
	var nilCh chan int
	select {
	case <-nilCh:
		fmt.Println("read")
	case nilCh <- 1:
		fmt.Println("write")
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}

// 非阻塞
// 通过使用select的default分支配合管道，我们可以实现非阻塞的收发操作，如下所示
func TrySend(ch chan int, ele int) bool {
	select {
	case ch <- ele:
		return true
	default:
		return false
	}
}
func TryRecv(ch chan int) (int, bool) {
	select {
	case v, ok := <-ch:
		return v, ok
	default:
		return 0, false
	}
}

// 同理，也可以实现非阻塞的判断一个context是否已经结束
func IsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// ==============锁
// Go 中 sync 包下的 Mutex 与 RWMutex 提供了互斥锁与读写锁两种实现，且提供了非常简单易用的 API，加锁只需要 Lock()，解锁也只需要 Unlock()。
// 需要注意的是，Go 所提供的锁都是非递归锁，也就是不可重入锁，所以重复加锁或重复解锁都会导致 fatal。锁的意义在于保护不变量，加锁是希望数据不会被其他协程修改
// 递归锁的话  重复加锁解锁都会导致死锁

// 互斥锁
// sync.Mutex 是 Go 提供的互斥锁实现，其实现了 sync.Locker

func Test17() {
	var wait sync.WaitGroup
	var num = 0

	var lock sync.Mutex
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func(data *int) {
			//枷锁
			lock.Lock()

			//模拟访问耗时
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			// 访问数据
			temp := *data
			// 模拟计算耗时
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			ans := 1
			//修改数据
			*data = temp + ans
			fmt.Println(*data)
			lock.Unlock()
			wait.Done()
		}(&num)
	}
	wait.Wait()
	fmt.Println("最终结果", num)
}

// 读写锁
// 互斥锁适合读操作与写操作频率都差不多的情况，对于一些读多写少的数据，如果使用互斥锁，会造成大量的不必要的协程竞争锁，这会消耗很多的系统资源，这时候就需要用到读写锁，即读写互斥锁，对于一个协程而言：
// 如果获得了读锁，其他协程进行写操作时会阻塞，其他协程进行读操作时不会阻塞
// 如果获得了写锁，其他协程进行写操作时会阻塞，其他协程进行读操作时会阻塞
// Go 中读写互斥锁的实现是 sync.RWMutex，它也同样实现了 Locker 接口
// // 加读锁
// func (rw *RWMutex) RLock()
// // 尝试加读锁
// func (rw *RWMutex) TryRLock() bool
// // 解读锁
// func (rw *RWMutex) RUnlock()
// // 加写锁
// func (rw *RWMutex) Lock()
// // 尝试加写锁
// func (rw *RWMutex) TryLock() bool
// // 解写锁
// func (rw *RWMutex) Unlock()
// 其中 TryRlock 与 TryLock 两个尝试加锁的操作是非阻塞式的，成功加锁会返回 true，无法获得锁时并不会阻塞而是返回 false。
// 读写互斥锁内部实现依旧是互斥锁，并不是说分读锁和写锁就有两个锁，从始至终都只有一个锁

var wait18 sync.WaitGroup
var count18 = 0
var rw sync.RWMutex

var cond = sync.NewCond(rw.RLocker())

func Test18() {
	wait18.Add(12)
	//读多写少
	go func() {
		for i := 0; i < 3; i++ {
			go Write(&count18)
		}
		wait18.Done()
	}()
	go func() {
		for i := 0; i < 7; i++ {
			go Read(&count18)
		}
		wait18.Done()
	}()
	wait18.Wait()
	fmt.Println("最终结果", count18)
}

func Read(i *int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	rw.RLock()
	fmt.Println("拿到读锁")
	for *i < 3 {
		cond.Wait()
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	fmt.Println("释放读锁", *i)
	rw.RUnlock()
	wait18.Done()
}
func Write(i *int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	rw.Lock()
	fmt.Println("拿到写锁")
	temp := *i
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	*i = temp + 1
	fmt.Println("释放写锁", *i)
	rw.Unlock()
	cond.Broadcast() //唤醒所有因条件变量阻塞的协程
	wait18.Done()
}

//对于锁而言，不应该将其作为值传递和存储，应该使用指针。

//条件变量
//条件变量，与互斥锁一同出现和使用，所以有些人可能会误称为条件锁，但它并不是锁，是一种通讯机制。Go 中的 sync.Cond 对此提供了实现
//func NewCond(l Locker) *Cond
//可以看到创建一个条件变量前提就是需要创建一个锁
//对于条件变量，应该使用 for 而不是 if，应该使用循环来判断条件是否满足，因为协程被唤醒时并不能保证当前条件就已经满足了。

// ======sync
// Go 中很大一部分的并发相关的工具都是 sync 标准库提供的，上述已经介绍过了 sync.WaitGroup，sync.Locker 等，除此之外，sync 包下还有一些其他的工具可以使用
// Once
// 当在使用一些数据结构时，如果这些数据结构太过庞大，可以考虑采用懒加载的方式，即真正要用到它的时候才会初始化该数据结构
// sync.Once Once 译为一次，sync.Once 保证了在并发条件下指定操作只会执行一次。它的使用非常简单，只对外暴露了一个 Do 方法
type MySlice struct {
	s []int
	o sync.Once
}

func (m *MySlice) Get(i int) (int, bool) {
	if m.s == nil {
		return 0, false
	} else {
		return (m.s)[i], true
	}
}
func (m *MySlice) Add(i int) {
	m.o.Do(func() {
		fmt.Println("初始化")
		if m.s == nil {
			m.s = make([]int, 0, 10)
		}
	})

	m.s = append(m.s, i)
}
func (m *MySlice) Len() int {
	return len(m.s)
}

func Test19() {
	var wait sync.WaitGroup
	var slice MySlice
	wait.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			slice.Add(1)
			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(slice.Len())
}

// pool
// sync.Pool 的设计目的是用于存储临时对象以便后续的复用，是一个临时的并发安全对象池，将暂时用不到的对象放入池中，
// 在后续使用中就不需要再额外的创建对象可以直接复用，减少内存的分配与释放频率，最重要的一点就是降低 GC 压力
// // 申请一个对象
// func (p *Pool) Get() any
// // 放入一个对象
// func (p *Pool) Put(x any)
// 并且 sync.Pool 有一个对外暴露的 New 字段，用于对象池在申请不到对象时初始化一个对象
// New func() any
type BigMemDate struct { //// BigMemData 假设这是一个占用内存很大的结构体
	M string
}

func Test20() {
	var wait sync.WaitGroup
	var pool sync.Pool
	var numOfObject atomic.Int64

	pool.New = func() any {
		numOfObject.Add(1)
		return BigMemDate{"大内存"}
	}
	wait.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			val := pool.Get()
			_ = val.(BigMemDate)
			pool.Put(val)
			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(numOfObject.Load())
}

//在使用 sync.Pool 时需要注意几个点：
//临时对象：sync.Pool 只适合存放临时对象，池中的对象可能会在没有任何通知的情况下被 GC 移除，所以并不建议将网络链接，数据库连接这类存入 sync.Pool 中。
//不可预知：sync.Pool 在申请对象时，无法预知这个对象是新创建的还是复用的，也无法知晓池中有几个对象
//并发安全：官方保证 sync.Pool 一定是并发安全，但并不保证用于创建对象的 New 函数就一定是并发安全的，New 函数是由使用者传入的，所以 New 函数的并发安全性要由使用者自己来维护，这也是为什么上例中对象计数要用到原子值的原因。
//最后需要注意的是，当使用完对象后，一定要释放回池中，如果用了不释放那么对象池的使用将毫无意义。

// map
// sync.Map 是官方提供的一种并发安全 Map 的实现，开箱即用
// // 根据一个key读取值，返回值会返回对应的值和该值是否存在
// func (m *Map) Load(key any) (value any, ok bool)
// // 存储一个键值对
// func (m *Map) Store(key, value any)
// // 删除一个键值对
// func (m *Map) Delete(key any)
// // 如果该key已存在，就返回原有的值，否则将新的值存入并返回，当成功读取到值时，loaded为true，否则为false
// func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool)
// // 删除一个键值对，并返回其原有的值，loaded的值取决于key是否存在
// func (m *Map) LoadAndDelete(key any) (value any, loaded bool)
// // 遍历Map，当f()返回false时，就会停止遍历
// func (m *Map) Range(f func(key, value any) bool)
func Test21() {
	var syncMap sync.Map
	syncMap.Store("a", 1)
	syncMap.Store("a", "a")
	fmt.Println(syncMap.Load("a"))                       // 读取数据
	fmt.Println(syncMap.LoadAndDelete("a"))              // 读取并删除
	fmt.Println(syncMap.LoadOrStore("a", "hello world")) // 读取或存入
	syncMap.Store("b", 2)
	//// 遍历map
	syncMap.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

// 并发
func Test22() {
	//下列操作大概率会异常 ：fatal error: concurrent map writes
	//myMap := make(map[int]int, 10)
	//var wait sync.WaitGroup
	//wait.Add(10)
	//for i := 0; i < 10; i++ {
	//	go func(n int) {
	//		for i := 0; i < 100; i++ {
	//			myMap[n] = n
	//		}
	//		wait.Done()
	//	}(i)
	//}
	//wait.Wait()
	var wait sync.WaitGroup
	var syncMap sync.Map
	wait.Add(10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			for i := 0; i < 100; i++ {
				syncMap.Store(n, n)
			}
			wait.Done()
		}(i)
	}
	wait.Wait()
	syncMap.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

// ====原子
// 类型 Go 标准库 sync/atomic 包下已经提供了原子操作相关的 API
//atomic.Bool{} atomic.Pointer[]{} atomic.Int32{} atomic.Int64{} atomic.Uint32{} atomic.Uint64{} atomic.Uintptr{} atomic.Value{}
//atmoic 包下原子操作只有函数签名，没有具体实现，具体的实现是由 plan9 汇编编写。

// 使用
// 每一个原子类型都会提供以下三个方法：
// Load()：原子的获取值
// Swap(newVal type) (old type)：原子的交换值，并且返回旧值
// Store(val type)：原子的存储值
// 不同的类型可能还会有其他的额外方法，比如整型类型都会提供 Add 方法来实现原子加减操作
func Test23() {
	var aint64 atomic.Uint64
	aint64.Store(64)
	aint64.Swap(128)
	aint64.Add(2)
	fmt.Println(aint64.Load())

	var aint32 int32
	atomic.StoreInt32(&aint32, 32)
	atomic.SwapInt32(&aint32, 64)
	atomic.AddInt32(&aint32, 1)
	fmt.Println(atomic.LoadInt32(&aint32))
}

// CAS
// atomic 包还提供了 CompareAndSwap 操作，也就是 CAS。它是实现乐观锁和无锁数据结构的核心。
// 乐观锁本身并不是锁，是一种并发条件下无锁化并发控制方式：线程/协程在修改数据前，不会先加锁，而是先读取数据，进行计算，
// 然后在提交修改时使用CAS来判断在此期间是否有其他线程修改过该数据。如果没有（值仍等于之前读取的值），则修改成功；否则，失败并重试。
// 因此之所以被称作乐观锁，是因为它总是乐观的假设共享数据不会被修改，仅当发现数据未被修改时才会去执行对应操作，而前面了解到的互斥量就是悲观锁，
// 互斥量总是悲观的认为共享数据肯定会被修改，所以在操作时会加锁，操作完毕后就会解锁。由于无锁化实现的并发，其安全性和效率相对于锁要高一些，许多并发安全的数据结构都采用了 CAS 来进行实现
var num int

func Add1(n int) {
	var lock sync.Mutex
	lock.Lock()
	num += n
	lock.Unlock()
}

var num64 int64

func Add2(n int64) {
	for {
		expect := atomic.LoadInt64(&num64)
		if atomic.CompareAndSwapInt64(&num64, expect, expect+n) {
			break
		}
	}
}

//大多数情况下，仅仅只是比较值是无法做到并发安全的，比如因 CAS 引起 ABA 问题，就需要使用额外加入 version 来解决问题。

// value
// atomic.Value 结构体，可以存储任意类型的值 但是它不能存储 nil，并且前后存储的值类型应当一致
func Test24() {
	var val atomic.Value
	val.Store(nil)
	fmt.Println(val.Load())
	// panic: sync/atomic: store of nil value into Value

	val.Store("hello world")
	val.Store(42)
	fmt.Println(val.Load())
	// panic: sync/atomic: store of inconsistently typed value into Value
}
