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