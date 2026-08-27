'''
import psutil

print("--- 💡 Cloud-FinOps Engine: 자원 모니터링 시작 ---")

# CPU와 Ram 사용량 읽어오기
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"현재 CPU 사용량: {cpu_usage}%")
print(f"현재 메모리 사용량: {memory_info.percent}%")
'''
import psutil

print("--- 💡 Cloud-FinOps Engine: 자원 모니터링 및 방치 프로세스 감지 ---")

# 1. 전체 자원 사용량
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"전체 CPU 사용량: {cpu_usage}%")
print(f"전체 메모리 사용량: {memory_info.percent}%\n")

# 2. 메모리를 많이 차지하는 상위 3개 프로그램 추적
print("🔍 [상위 메모리 점유 프로세스 TOP 3]")
processes = []

for proc in psutil.process_iter(['pid', 'name', 'memory_percent']):
    try:
        processes.append(proc.info)
    except (psutil.NoSuchProcess, psutil.AccessDenied):
        pass

# 메모리 사용량 순으로 정렬
top_processes = sorted(processes, key=lambda x: x['memory_percent'], reverse=True)[:3]

for proc in top_processes:
    print(f"📌 PID: {proc['pid']} | 이름: {proc['name']} | 메모리 점유율: {proc['memory_percent']:.2f}%")
    