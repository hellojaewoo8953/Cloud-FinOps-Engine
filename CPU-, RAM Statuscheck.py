import psutil
import time

print("--- 💡 Cloud-FinOps Engine: 자원 모니터링 시작 ---")

# 현재 CPU 및 메모리 사용량 읽어오기
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"현재 CPU 사용량: {cpu_usage}%")
print(f"현재 메모리 사용량: {memory_info.percent}%")

# 메모리가 80% 이상이거나 CPU가 높은 상태인지 간단 점검
if cpu_usage > 80:
    print("⚠️ 경고: CPU 사용량이 너무 높습니다!")
else:
    print("✅ CPU 상태가 안정적입니다.")
    