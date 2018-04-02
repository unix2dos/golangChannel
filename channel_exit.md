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

