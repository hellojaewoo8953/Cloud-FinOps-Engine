FROM python:3.10-slim

WORKDIR /app

# 보안 업데이트 및 필수 라이브러리 설치
RUN apt-get update && apt-get install -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && pip install --no-cache-dir psutil requests python-dotenv

# 비권한 사용자 생성 및 적용
RUN useradd -m appuser
COPY --chown=appuser:appuser . /app
USER appuser

CMD ["python", "Prozessdaten sammeln.py"]

