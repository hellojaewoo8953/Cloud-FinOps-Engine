package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1025")
	if err != nil {
		fmt.Println("서버 실행 실패:", err)
		return
	}
	defer listener.Close()

	fmt.Println("📧 [가짜 SMTP 서버] 1025번 포트에서 메일 수신 대기 중...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// 접속 응답
	writer.WriteString("220 Fake SMTP Server Ready\r\n")
	writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		fmt.Println("📩 [수신 내용]:", line)

		if strings.HasPrefix(line, "HELO") || strings.HasPrefix(line, "EHLO") {
			writer.WriteString("250 Hello\r\n")
		} else if strings.HasPrefix(line, "MAIL FROM:") {
			writer.WriteString("250 OK\r\n")
		} else if strings.HasPrefix(line, "RCPT TO:") {
			writer.WriteString("250 OK\r\n")
		} else if line == "DATA" {
			writer.WriteString("354 Start mail input\r\n")
		} else if line == "." {
			writer.WriteString("250 OK Message accepted\r\n")
		} else if line == "QUIT" {
			writer.WriteString("221 Bye\r\n")
			writer.Flush()
			return
		}
		writer.Flush()
	}
}