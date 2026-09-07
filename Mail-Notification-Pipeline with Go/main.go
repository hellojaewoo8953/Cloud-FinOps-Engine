package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func alertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🚨 [알림] 파이썬 감시원으로부터 자원 낭비 경고 수신!")

	conn, err := net.Dial("tcp", "127.0.0.1:1025")
	if err != nil {
		fmt.Println("❌ 메일 전송 실패 (우체통 연결 불가):", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to send email alert"))
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 서버 220 대기
	reader.ReadString('\n')

	commands := []string{
		"HELO localhost\r\n",
		"MAIL FROM:<finops-engine@local.com>\r\n",
		"RCPT TO:<admin@company.com>\r\n",
		"DATA\r\n",
		"Subject: [FinOps Engine] Resource Waste Alert!\r\n\r\nWarning: CPU usage high.\r\n.\r\n",
		"QUIT\r\n",
	}

	for _, cmd := range commands {
		writer.WriteString(cmd)
		writer.Flush()
		reader.ReadString('\n')
	}

	fmt.Println("📧 [성공] 가짜 메일 서버로 알림 메일을 성공적으로 발송했습니다!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received and Email sent!"))
}

func main() {
	http.HandleFunc("/alert", alertHandler)
	fmt.Println("🚀 Go FinOps 조종기 서버가 8080번 포트에서 실행 중입니다...")
	http.ListenAndServe(":8080", nil)
}