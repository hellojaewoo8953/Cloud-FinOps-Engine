package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

func sanitize(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func main() {
	smtpHost := getEnvOrDefault("SMTP_HOST", "127.0.0.1")
	smtpPort := getEnvOrDefault("SMTP_PORT", "1025")

	from := sanitize(getEnvOrDefault("SMTP_FROM", "alert@example.com"))
	toEmail := sanitize(getEnvOrDefault("SMTP_TO", "admin@example.com"))

	if from == "" || toEmail == "" {
		log.Fatal("❌ 이메일 주소 설정이 올바르지 않습니다.")
	}

	to := []string{toEmail}

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: [경고] CPU 사용량 초과!\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
		"서버 자원이 낭비되고 있습니다.", from, toEmail))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, nil, from, to, msg)
	if err != nil {
		log.Fatalf("❌ 메일 전송 실패: %v", err)
	}
	fmt.Println("✅ 성공적으로 메일 서버에 메일을 보냈습니다!")
}
