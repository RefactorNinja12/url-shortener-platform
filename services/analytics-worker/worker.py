import json
import logging
import os
import time

import pika
import psycopg2

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger("analytics-worker")

QUEUE_NAME = "click_events"


def parse_click_event(body: bytes) -> tuple[str, str, str]:
    """Parsear ett click_events-meddelande till (code, original_url, clicked_at)."""
    data = json.loads(body)
    return data["code"], data["original_url"], data["clicked_at"]


def insert_click_event(conn, code: str, original_url: str, clicked_at: str) -> None:
    with conn.cursor() as cur:
        cur.execute(
            "INSERT INTO click_events (code, original_url, clicked_at) VALUES (%s, %s, %s)",
            (code, original_url, clicked_at),
        )
    conn.commit()


def connect_db(database_url: str, retries: int = 10, delay: float = 2.0):
    last_err: Exception | None = None
    for _ in range(retries):
        try:
            return psycopg2.connect(database_url)
        except psycopg2.OperationalError as err:
            last_err = err
            logger.warning("database not ready yet, retrying: %s", err)
            time.sleep(delay)
    raise last_err


def connect_rabbitmq(rabbitmq_url: str, retries: int = 10, delay: float = 2.0):
    params = pika.URLParameters(rabbitmq_url)
    last_err: Exception | None = None
    for _ in range(retries):
        try:
            return pika.BlockingConnection(params)
        except pika.exceptions.AMQPConnectionError as err:
            last_err = err
            logger.warning("rabbitmq not ready yet, retrying: %s", err)
            time.sleep(delay)
    raise last_err


def main() -> None:
    database_url = os.environ.get(
        "DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/urlshortener"
    )
    rabbitmq_url = os.environ.get("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

    conn = connect_db(database_url)
    rabbit_conn = connect_rabbitmq(rabbitmq_url)
    channel = rabbit_conn.channel()
    channel.queue_declare(queue=QUEUE_NAME, durable=True)

    def on_message(ch, method, _properties, body):
        try:
            code, original_url, clicked_at = parse_click_event(body)
            insert_click_event(conn, code, original_url, clicked_at)
            logger.info("stored click event for code=%s", code)
            ch.basic_ack(delivery_tag=method.delivery_tag)
        except Exception:
            logger.exception("failed to process click event, requeueing")
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

    channel.basic_qos(prefetch_count=10)
    channel.basic_consume(queue=QUEUE_NAME, on_message_callback=on_message)

    logger.info("analytics-worker started, waiting for click events")
    channel.start_consuming()


if __name__ == "__main__":
    main()
