package main

import (
	"GoRedis/tcp"
)

var cfg tcp.Config = tcp.Config{
	Address: ":8000",
}

func main() {
	echoHandler := tcp.NewEchoHandler()
	tcp.ListenAndServeWithSignal(&cfg, echoHandler)
}

// func ListenAndServe(address string) {
// 	// 绑定监听地址
// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		log.Fatalf("listen err: %v", err)
// 	}
// 	defer listener.Close()
// 	log.Printf("bind: %s, start listening...", address)

// 	for {
// 		// Accept会一直阻塞到有新的连接或者连接中断
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			// 通常是listener被关闭，无法继续监听
// 			log.Fatalf("accept err : %v", err)
// 		}
// 		go Handle(conn)
// 	}
// }

// func Handle(conn net.Conn) {
// 	// bufio提供缓冲区
// 	reader := bufio.NewReader(conn)
// 	for {
// 		// ReadString 会一直阻塞直到遇到分隔符 '\n'
// 		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
// 		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
// 		msg, err := reader.ReadString('\n')
// 		if err != nil {
// 			// 连接中断或者是关闭
// 			if err == io.EOF {
// 				log.Println("connnection close")
// 			} else {
// 				log.Println(err)
// 			}
// 			return
// 		}
// 		b := []byte(msg)
// 		conn.Write(b)
// 	}
// }
