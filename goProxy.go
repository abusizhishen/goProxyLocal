package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// 定义命令行参数
	localAddr := flag.String("local", "", "Local address to proxy (format: ip:port)")
	externalAddr := flag.String("external", "0.0.0.0:9090", "External address to bind (format: ip:port)")
	flag.Parse()

	// 校验参数
	if *localAddr == "" {
		fmt.Println("Error: -local argument is required")
		flag.Usage()
		os.Exit(1)
	}

	// 启动代理服务
	go startProxy(*localAddr, *externalAddr)

	log.Printf("Proxying from %s to %s\n", *localAddr, *externalAddr)

	select {} // 阻止程序退出
}

func startProxy(localAddress, externalAddress string) {
	// 监听外部地址
	listener, err := net.Listen("tcp", externalAddress)
	if err != nil {
		log.Fatalf("Failed to bind on %s: %v", externalAddress, err)
	}
	defer listener.Close()

	log.Printf("Listening on %s\n", externalAddress)

	for {
		// 接收外部连接
		externalConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 连接到本地地址
		localConn, err := net.Dial("tcp", localAddress)
		if err != nil {
			log.Printf("Failed to connect to local address %s: %v", localAddress, err)
			externalConn.Close()
			continue
		}

		// 转发数据
		go handleConnection(externalConn, localConn)
		go handleConnection(localConn, externalConn)
	}
}

func handleConnection(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()

	_, err := io.Copy(dest, src)
	if err != nil {
		log.Printf("Data transfer error: %v", err)
	}
}
