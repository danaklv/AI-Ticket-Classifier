import json
import uuid
import time
from kafka.consumer import KafkaConsumer
from kafka.producer import KafkaProducer
from classifier.classifier import predict_category
from utils.logger import setup_logger

logger = setup_logger("main")

BROKER = "localhost:9094"
INPUT_TOPIC = "tickets"
OUTPUT_TOPIC = "classified_tickets"
GROUP_ID = "python-classifier-group"

consumer = KafkaConsumer(BROKER, INPUT_TOPIC, GROUP_ID)
producer = KafkaProducer(BROKER)

logger.info("Waiting for tickets...")

try:
    while True:
        msg = consumer.poll_message()
        if not msg:
            continue

        ticket = json.loads(msg.value().decode("utf-8"))
        text = ticket.get("text", "")
        logger.info("Received ticket %s: %s", ticket.get("id"), text[:40])

        category = predict_category(text, text)
        result = {
            "event_id": str(uuid.uuid4()),
            "ticket_id": ticket.get("id"),
            "text": text,
            "predicted_category": category,
            "processed_at": time.strftime("%Y-%m-%d %H:%M:%S"),
        }

        producer.send(OUTPUT_TOPIC, result)

except KeyboardInterrupt:
    logger.info("Stopping consumer...")
finally:
    consumer.close()
    producer.flush()
