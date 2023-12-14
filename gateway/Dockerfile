FROM python:3.11-bookworm

WORKDIR /gateway

COPY ./requirements.txt .

RUN pip install --no-cache-dir --upgrade -r ./requirements.txt

COPY . .

EXPOSE 9080

ENV FLASK_APP=gateway.py

CMD ["python", "-m", "flask", "run", "--host=0.0.0.0", "--port=9080", "--debug"]