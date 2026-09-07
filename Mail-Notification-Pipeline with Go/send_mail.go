package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	// 보낸 사람, 받는 사람, 메일 내용 작성
	from := "alert@finops.com"
	to := []string{"admin@company.com"}
	message := []byte("Subject: [경고] CPU 사용량 초과!\r\n\r\n서버 자원이 낭비되고 있습니다.")

	// 우리가 만든 1025번 가짜 우체통으로 메일 발송!
	err := smtp.SendMail("127.0.0.1:1025", nil, from, to, message)
	if err != nil {
		fmt.Println("❌ 메일 전송 실패:", err)
		return
	}
	fmt.Println("✅ 성공적으로 가짜 메일 서버에 메일을 보냈습니다!")
}
