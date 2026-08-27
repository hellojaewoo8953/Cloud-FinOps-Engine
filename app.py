'''
import psutil

print("--- 💡 Cloud-FinOps Engine: Resource Monitoring ---")

# Overall Resource Usage
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"Total CPU Usage: {cpu_usage}%")
print(f"Total Memory Usage: {memory_info.percent}%")
'''



import psutil

print("--- 💡 Cloud-FinOps Engine: Resource Monitoring & Idle Process Detection ---")

# 1. Overall Resource Usage
cpu_usage = psutil.cpu_percent(interval=1)
memory_info = psutil.virtual_memory()

print(f"Total CPU Usage: {cpu_usage}%")
print(f"Total Memory Usage: {memory_info.percent}%\n")

# 2. Track Top 3 Memory-Consuming Processes
print("🔍 [Top 3 Memory-Consuming Processes]")
processes = []

for proc in psutil.process_iter(['pid', 'name', 'memory_percent']):
    try:
        processes.append(proc.info)
    except (psutil.NoSuchProcess, psutil.AccessDenied):
        pass

# Sort by memory usage
top_processes = sorted(processes, key=lambda x: x['memory_percent'], reverse=True)[:3]

for proc in top_processes:
    print(f"📌 PID: {proc['pid']} | Name: {proc['name']} | Memory Usage: {proc['memory_percent']:.2f}%")