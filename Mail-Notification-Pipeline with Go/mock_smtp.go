package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	maxConcurrentConns  = 20
	maxMessageSizeBytes = 1024 * 1024 // 1MB
	maxLineSizeBytes    = 4096        // 4KB
	readTimeout         = 10 * time.Second
)

var sem = make(chan struct{}, maxConcurrentConns)

func main() {
	port := "1025"
	if envPort := os.Getenv("SMTP_SERVER_PORT"); envPort != "" {
		port = envPort
	}

	// 127.0.0.1에 강제 바인딩 (외부 노출 절대 차단)
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		fmt.Println("서버 실행 실패:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("📧 [Mock SMTP 서버] 127.0.0.1:%s 대기 중...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		select {
		case sem <- struct{}{}:
			go func(c net.Conn) {
				defer func() { <-sem }()
				handleClient(c)
			}(conn)
		default:
			conn.Close()
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(readTimeout))

	limitedReader := io.LimitReader(conn, maxMessageSizeBytes)
	reader := bufio.NewReader(limitedReader)
	writer := bufio.NewWriter(conn)

	writer.WriteString("220 Fake SMTP Server Ready\r\n")
	writer.Flush()

	inDataMode := false
	totalDataBytes := 0

	for {
		_ = conn.SetReadDeadline(time.Now().Add(readTimeout))

		lineBytes, isPrefix, err := reader.ReadLine()
		if err != nil {
			return
		}

		if isPrefix {
			writer.WriteString("500 Line too long\r\n")
			writer.Flush()
			return
		}

		trimmedLine := strings.TrimSpace(string(lineBytes))

		if inDataMode {
			totalDataBytes += len(lineBytes)
			// DATA 모드 중 최대 허용 용량 초과 시 즉시 차단
			if totalDataBytes > maxMessageSizeBytes {
				writer.WriteString("552 Message size exceeds limit\r\n")
				writer.Flush()
				return
			}

			if trimmedLine == "." {
				inDataMode = false
				writer.WriteString("250 OK Message accepted\r\n")
				writer.Flush()
			}
			continue
		}

		upperLine := strings.ToUpper(trimmedLine)

		if strings.HasPrefix(upperLine, "HELO") || strings.HasPrefix(upperLine, "EHLO") {
			writer.WriteString("250 Hello\r\n")
		} else if strings.HasPrefix(upperLine, "MAIL FROM:") {
			writer.WriteString("250 OK\r\n")
		} else if strings.HasPrefix(upperLine, "RCPT TO:") {
			writer.WriteString("250 OK\r\n")
		} else if upperLine == "DATA" {
			inDataMode = true
			writer.WriteString("354 Start mail input; end with <CRLF>.<CRLF>\r\n")
		} else if upperLine == "QUIT" {
			writer.WriteString("221 Bye\r\n")
			writer.Flush()
			return
		} else {
			writer.WriteString("500 Command unrecognized\r\n")
		}
		writer.Flush()
	}
}