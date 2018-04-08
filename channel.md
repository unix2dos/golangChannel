### 创建一个tcp

因为创建的是tcp,用http有点问题, 用nc测试

```
func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		_, err := io.WriteString(conn, time.Now().Format("2006-01-02 15:04:05\n"))
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(time.Second)
	}
}

func main() {

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		handleConn(conn)
	}
}
```



连接tcp

```
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	copy(os.Stdout, conn)
}

func copy(dst io.Writer, src io.Reader) {

	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
```

### golang 闭包循环变量快照问题

##### 闭包都是匿名函数
##### goroutine defer fuc指针


+ 匿名函数存储

```
func test() []func() {
    var s []func()

    for i := 0; i < 3; i++ {
        s = append(s, func() {
            fmt.Println(&i, i)
        })
    }

    return s
}
func main() {
    for _, f := range test() {
        f()   //全部输出3
    }
}
```

+ goroutine

```
	s := []string{"a", "b", "c"}
	for _, v := range s {
		go func() {
			fmt.Println(v)
		}()
	}
	// for range 引用, 全部输出c
```

+ defer调用

```
func main() {
    x, y := 1, 2

    defer func(a int) { 
        fmt.Printf("x:%d,y:%d\n", a, y)
    }(x)

    x += 100
    y += 100
    fmt.Println(x, y) //输出101, 102  x:1,y:102
}
```



### channel


+ 接收者先收到, 再唤醒发送者

	接收者收到数据发生在唤醒发送者goroutine之前



### 数据竞争

+ 避免数据竞争
	+ 全局变量一开始就初始化, 其他goroutine都读取
	+ 避免多个goroutine访问同一个变量, 变量限定一个单独的goroutine

		> 由于其它的goroutine不能够直接访问变量，它们只能使用一个channel来发送给指定的goroutine请求来查询更新变量。
		
		> 这也就是Go的口头禅“不要使用共享数据来通信；使用通信来共享数据”。
	+ 避免数据竞争的方法是允许很多goroutine去访问变量，但是在同一个时刻最多只有一个goroutine在访问。这种方式被称为“互斥”

	
	
	
### 内存同步	
	
```
var x, y int
go func() {
    x = 1 // A1
    fmt.Print("y:", y, " ") // A2
}()
go func() {
    y = 1                   // B1
    fmt.Print("x:", x, " ") // B2
}()


x:0 y:0
y:0 x:0
```

因为赋值和打印指向不同的变量，编译器可能会断定两条语句的顺序不会影响执行结果，并且会交换两个语句的执行顺序。

如果两个goroutine在不同的CPU上执行，每一个核心有自己的缓存，这样一个goroutine的写入对于其它goroutine的Print，在主存同步之前就是不可见的了。


### 竞争检测

只要在go build，go run或者go test命令后面加上-race的flag，就会使编译器创建一个你的应用的“修改”版

```
go run -race cache/memo1_goroutine.go
```

### 线程和gorutine


+ 线程一般2MB
+ gorutine开始2KB 最大值1GB
+ n个线程调度m个gorutine
+ GOMAXPROCS 是多少个线程执行Go的代码





### 优雅的等待gorutine退出, 就是死等

+ 阻塞等待一个goroutine
+ 多个使用 waitgroup


### 圣经里是并发的退出goroutine干的事情, 通知退出
	
+ 干事的gorutine for select case 专门发送结束的channle, 获取后清空一些资源
+ 其他gorutine就check专门发送结束的channel return
+ 这样gorutine会迅速停止


----

+ 在不能更改channel状态的情况下，没有简单普遍的方式来检查channel是否已经关闭了

+ 关闭已经关闭的channel会导致panic，所以在closer(关闭者)不知道channel是否已经关闭的情况下去关闭channel是很危险的

+ 发送值到已经关闭的channel会导致panic，所以如果sender(发送者)在不知道channel是否已经关闭的情况下去向channel发送值是很危险的


解决方法:

+ _,ok := <- jobs  测试jobs是否关闭, 这种方法不对

此时如果 channel 关闭，ok 值为 false，如果 channel 没有关闭，则会漏掉一个 jobs

+ 使用select

新建一个channle,达到条件发送到这个channel. 干活的gorutine 通过 select这个channel然后干结束的事情


### 优雅的关闭channel, 是为了不崩溃
> Channel关闭原则

+ 不要在消费端关闭channel (生产方不知道继续发)

+ 不要在有多个并行的生产者时对channel执行关闭操作。(其他生成不知道继续发)

+ 也就是说应该只在[唯一的或者最后唯一剩下]的生产者协程中关闭channel，来通知消费者已经没有值可以继续读了。只要坚持这个原则，就可以确保向一个已经关闭的channel发送数据的情况不可能发生。


> 暴力关闭channel的正确方法

```
func SafeClose(ch chan T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	
	// assume ch != nil here.
	close(ch) // panic if ch is closed
	return true // <=> justClosed = true; return
}
```

```
 func SafeSend(ch chan T, value T) (closed bool) {
	defer func() {
		if recover() != nil {
			// The return result can be altered 
			// in a defer function call.
			closed = true
		}
	}()
	
	ch <- value // panic if ch is closed
	return false // <=> closed = false; return
}
```

> 优雅的关闭channel的方法

上文的SafeSend方法一个很大的劣势在于它不能用在select块的case语句中。而另一个很重要的劣势在于像我这样对代码有洁癖的人来说，使用panic/recover和sync/mutex来搞定不是那么的优雅。下面我们引入在不同的场景下可以使用的纯粹的优雅的解决方法。

+ 多个消费者，单个生产者。这种情况最简单，直接让生产者关闭channel好了。

```
func main() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumCosters = 100

	dataCh := make(chan int, 100) //数据

	wg := sync.WaitGroup{}
	wg.Add(NumCosters)

	//produce
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// the only produce can close the channel safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	//cost
	for i := 0; i < NumCosters; i++ {
		go func() {
			defer wg.Done()

			for value := range dataCh {
				log.Println(value)
			}

		}()
	}

	wg.Wait()
}

```
+ 多个生产者，单个消费者。这种情况要比上面的复杂一点。我们不能在消费端关闭channel，因为这违背了channel关闭原则。但是我们可以让消费端关闭一个附加的信号来通知发送端停止生产数据。

```
func main() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumProduces = 100

	dataCh := make(chan int, 100) //数据
	stopCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(1)

	//produce
	for i := 0; i < NumProduces; i++ {
		go func() {
			for {
				select {
				case <-stopCh:
					return
				default:
					value := rand.Intn(MaxRandomNumber)
					dataCh <- value
				}

			}
		}()
	}

	//cost
	go func() {
		defer wg.Done()
		for value := range dataCh {
			if value == 0 {
				close(stopCh)
				return
			}

			log.Println(value)
		}

	}()

	wg.Wait()
}
```

生产者同时也是退出信号channel的接受者，退出信号channel仍然是由它的生产端关闭的，所以这仍然没有违背channel关闭原则。值得注意的是，`这个例子中生产端和接受端都没有关闭消息数据的channel，`channel在没有任何goroutine引用的时候会自行关闭，而不需要显示进行关闭。


+ 多个生产者，多个消费者

我们不能让任意的receivers和senders关闭data channel，也不能让任何一个receivers通过关闭一个额外的signal channel来通知所有的senders和receivers退出游戏。这么做的话会打破channel closing principle。但是，我们可以引入一个moderator来关闭一个额外的signal channel。

```
func main() {

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumProduces = 100
	const NumCosters = 10

	dataCh := make(chan int, 100) //数据
	stopCh := make(chan struct{})
	stopByNameCh := make(chan string, 1) //1是为了裁判goroutine没准备好就收到通知了
	var stopName string

	//裁判goroutine
	go func() {
		stopName = <-stopByNameCh
		close(stopCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(NumCosters)

	//produce
	for i := 0; i < NumProduces; i++ {
		go func(i int) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					select {
					case stopByNameCh <- strconv.Itoa(i): //加select是为了其他有发送的, 我就什么不干了
					default:
					}
					return
				}

				select { //这边写是为了尽早退出,越早越好
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(i)
	}

	//cost
	for i := 0; i < NumCosters; i++ {
		go func(i int) {
			defer wg.Done()

			for {

				select { //这边写是为了尽早
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						select {
						case stopByNameCh <- strconv.Itoa(i): //加select是为了不阻塞
						default:
						}
						return
					}
					log.Println(value)
				}
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("---> " + stopName)
}
```



### Context

通过context，我们可以方便地对同一个请求所产生地goroutine进行约束管理，可以设定超时、deadline，甚至是取消这个请求相关的所有goroutine。形象地说，假如一个请求过来，需要A去做事情，而A让B去做一些事情，B让C去做一些事情，A、B、C是三个有关联的goroutine，那么问题来了：假如在A、B、C还在处理事情的时候请求被取消了，那么该如何优雅地同时关闭goroutine A、B、C呢？这个时候就轮到context包上场了。



```
func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()

	fmt.Println(A(ctx))

	time.Sleep(time.Hour)
}

func A(ctx context.Context) string {

	go fmt.Println(B(ctx))

	for {
		select {
		case <-ctx.Done():
			return "A Done"
		}
	}
	return ""
}

func B(ctx context.Context) string {

	go fmt.Println(C(ctx))
	for {
		select {
		case <-ctx.Done():
			return "B Done"
		}
	}
	return ""
}

func C(ctx context.Context) string {
	for {
		select {
		case <-ctx.Done():
			return "C Done"
		}
	}
	return ""
}
```


这里的例子是直接调用了context.WithCancel()，我们也可以使用context.WithTimeout()和context.WithDeadline()来设置goroutine的超时时间和最终的运行时间。


另外有一个方法在例子中没有用到，那就是context.WithValue()。这个方法是用来传递在这次的请求处理中相关goroutine的共享变量，这与全局变量是有所区别的，因为它只在这次的请求范围内有效



> context的使用规范

+ 不要把context存储在结构体中，而是要显式地进行传递
+ 把context作为第一个参数，并且一般都把变量命名为ctx
+ 就算是程序允许，也不要传入一个nil的context，如果不知道是否要用context的话，用context.TODO()来替代
+ context.WithValue()只用来传递请求范围的值，不要用它来传递可选参数
+ 就算是被多个不同的goroutine使用，context也是安全的
