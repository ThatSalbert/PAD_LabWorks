FROM python:3.11-bookworm

WORKDIR /service-discovery

COPY ./requirements.txt .

RUN pip install --no-cache-dir --upgrade -r ./requirements.txt

COPY . .

EXPOSE 9000

ENV FLASK_APP=service_discovery.py

CMD ["python", "-m", "flask", "run", "--host=0.0.0.0", "--port=9000", "--debug"]