package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func sanitizeHeader(input string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(input)
}

// IP별 Rate Limiter 관리 구조체
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
	// 메모리 누수 방지: 오래된 IP 정리
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			limiter.mu.Lock()
			limiter.ips = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
		}
	}()
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}
	return limiter
}

var limiter = NewIPRateLimiter(rate.Every(1*time.Second), 3) // IP당 초당 1개, 최대 3개 버스트

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if !limiter.GetLimiter(ip).Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB 제한

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	apiToken := os.Getenv("ALERT_API_TOKEN")
	if apiToken == "" {
		log.Println("❌ [보안 경고] ALERT_API_TOKEN 환경변수가 설정되지 않았습니다.")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	expectedHeader := "Bearer " + apiToken

	if subtle.ConstantTimeCompare([]byte(authHeader), []byte(expectedHeader)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// [비동기 처리] 메일 발송 로직을 Goroutine으로 분리하여 HTTP 응답 지연(Blocking DoS) 방지
	go func() {
		smtpHost := getEnv("SMTP_HOST", "127.0.0.1")
		smtpPort := getEnv("SMTP_PORT", "1025")
		from := sanitizeHeader(getEnv("SMTP_FROM", "finops-engine@local.com"))
		to := sanitizeHeader(getEnv("SMTP_TO", "admin@company.com"))

		msg := []byte(fmt.Sprintf("To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: [FinOps Engine] Resource Waste Alert!\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
			"Warning: CPU usage high.\r\n", to, from))

		addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
		err := smtp.SendMail(addr, nil, from, []string{to}, msg)
		if err != nil {
			log.Printf("❌ 메일 비동기 전송 실패: %v\n", err)
			return
		}
		fmt.Println("📧 [성공] 가짜 메일 서버로 알림 메일 비동기 발송 완료!")
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received and processing email dispatch"))
}

func main() {
	port := getEnv("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/alert", rateLimitMiddleware(alertHandler))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Printf("🚀 Go FinOps 조종기 서버가 %s번 포트에서 실행 중입니다...\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("서버 실행 오류: %v", err)
	}
}