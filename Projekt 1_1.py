import psutil
import time

print("--- 💡 Cloud-FinOps Engine: Resource Monitoring & Idle Process Detection ---")

# 1. 전체 시스템 자원 모니터링
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"Total CPU Usage: {cpu_usage}%")
print(f"Total Memory Usage: {memory_info.percent}%")

# 2. 프로세스 정보 수집 (None 값 예외 처리)
processes = []
for proc in psutil.process_iter(['pid', 'name', 'cpu_percent', 'memory_percent']):
    try:
        pinfo = proc.info
        pinfo['memory_percent'] = pinfo['memory_percent'] or 0.0
        pinfo['cpu_percent'] = pinfo['cpu_percent'] or 0.0
        processes.append(pinfo)
    except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
        pass

# 3. 메모리 사용량이 가장 높은 상위 3개 프로세스 추출
top_processes = sorted(processes, key=lambda x: x['memory_percent'], reverse=True)[:3]

print("\n🔍 [Top 3 Memory-Consuming Processes]")
for p in top_processes:
    print(f"- PID: {p['pid']} | Name: {p['name']} | Memory: {p['memory_percent']:.2f}%")
    print(f"- PID: {p['pid']} | Name: {p['name']} | CPU: {p['cpu_percent']:.2f}%")