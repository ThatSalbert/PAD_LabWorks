FROM python:3.11-bookworm

WORKDIR /coordinator

COPY ./requirements.txt .

RUN pip install --no-cache-dir --upgrade -r ./requirements.txt

COPY . .

EXPOSE 9081

ENV FLASK_APP=coordinator.py

CMD ["python", "-m", "flask", "run", "--host=0.0.0.0", "--port=9081", "--debug"]