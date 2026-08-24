# analytics-worker

Python-worker som konsumerar klick-events från RabbitMQ (kön `click_events`,
publicerad av go-api vid varje redirect) och skriver dem till tabellen
`click_events` i Postgres.

## Köra lokalt

```bash
pip install -r requirements-dev.txt
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/urlshortener \
RABBITMQ_URL=amqp://guest:guest@localhost:5672/ \
python worker.py
```

## Tester

```bash
pip install -r requirements-dev.txt
pytest -v
```
