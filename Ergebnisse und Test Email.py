import os
import psutil
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.header import Header
from dotenv import load_dotenv
import time

# .env 파일에 있는 내용을 읽어옵니다.
load_dotenv()

# .env 파일에서 정보를 가져오고, 없으면 기본값을 씁니다.
SMTP_HOST = os.getenv('SMTP_HOST', 'localhost')
SMTP_PORT = int(os.getenv('SMTP_PORT', 1025))

# --------------------------------------------------
# 1. 가짜 메일 서버 테스트용 설정 (비밀번호 불필요!)
# --------------------------------------------------
SENDER_EMAIL = "test-bot@finops.local"      # 임의의 보내는 사람
RECEIVER_EMAIL = "developer@finops.local"    # 임의의 받는 사람

def send_email_alert(subject, body):
    """내 컴퓨터 안의 가짜 메일 서버(localhost:1025)로 이메일을 발송하는 함수"""
    msg = MIMEMultipart()
    msg['From'] = SENDER_EMAIL
    msg['To'] = RECEIVER_EMAIL
    msg['Subject'] = Header(subject, 'utf-8')

    # 한글 본문 인코딩 (utf-8 지정)
    msg.attach(MIMEText(body, 'plain', 'utf-8'))

    try:
        # 내 컴퓨터(localhost)의 1025번 포트에 떠 있는 가짜 메일 서버로 연결
        server = smtplib.SMTP(SMTP_HOST, SMTP_PORT)
        server.sendmail(SENDER_EMAIL, RECEIVER_EMAIL, msg.as_string())
        server.quit()
        
        print("📧 [가짜 메일 서버]로 이메일 알림 전송 성공!")
    except Exception as e:
        print(f"❌ 이메일 전송 실패: {e}")
        print("💡 팁: 가짜 메일 서버 터미널이 켜져 있는지 확인해 주세요!")

print("--- 💡 Cloud-FinOps Engine v2: 낭비 자원 감지 및 이메일 알림 ---")

# 2. 전체 시스템 자원 모니터링
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"Total CPU Usage: {cpu_usage}%")
print(f"Total Memory Usage: {memory_info.percent}%")

# 3. 낭비(Idle) 프로세스 감지 로직 (None 예외 처리 포함)
idle_processes = []

for proc in psutil.process_iter(['pid', 'name', 'cpu_percent', 'memory_percent']):
    try:
        pinfo = proc.info
        
        cpu_p = pinfo['cpu_percent'] if pinfo['cpu_percent'] is not None else 0.0
        mem_p = pinfo['memory_percent'] if pinfo['memory_percent'] is not None else 0.0
        
        if cpu_p < 1.0 and mem_p > 0.5:
            pinfo['cpu_percent'] = cpu_p
            pinfo['memory_percent'] = mem_p
            idle_processes.append(pinfo)
    except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
        pass

# 4. 결과 정리 및 이메일 발송
if idle_processes:
    top_idle = sorted(idle_processes, key=lambda x: x['memory_percent'], reverse=True)[:3]
    
    email_subject = f"⚠️ [Cloud-FinOps] 낭비 자원 감지 리포트 (메모리: {memory_info.percent}%)"
    
    email_body = f"안녕하세요,\n\nCloud-FinOps 자원 절감 엔진에서 시스템 이상 및 낭비 자원을 감지했습니다.\n\n"
    email_body += f"🖥 전체 메모리 사용량: {memory_info.percent}%\n"
    email_body += f"🖥 전체 CPU 사용량: {cpu_usage}%\n\n"
    email_body += f"🔍 방치된 프로세스 목록 (상위 3개):\n"
    
    for p in top_idle:
        email_body += f"- 프로그램 이름: {p['name']} | PID: {p['pid']} | 메모리 점유: {p['memory_percent']:.2f}%\n"
        
    email_body += f"\n확인 후 불필요한 프로세스를 정리해 주세요.\n감사합니다."
    
    print("\n--- 생성된 이메일 본문 ---")
    print(email_body)
    
    # 이메일 전송 함수 호출
    send_email_alert(email_subject, email_body)
else:
    print("✅ 현재 방치된 낭비 프로세스가 없습니다.")