from confluent_kafka import Consumer, Producer
from classifier import predict_category
import json
import time

KAFKA_BROKER = "localhost:9092"
INPUT_TOPIC = "tickets"
OUTPUT_TOPIC = "classified_tickets"

# Конфигурация Kafka
consumer_conf = {
    "bootstrap.servers": KAFKA_BROKER,
    "group.id": "python-classifier-group",
    "auto.offset.reset": "earliest"
}
producer_conf = {"bootstrap.servers": KAFKA_BROKER}

consumer = Consumer(consumer_conf)
producer = Producer(producer_conf)
consumer.subscribe([INPUT_TOPIC])

print("📡 Waiting for tickets...")

while True:
    msg = consumer.poll(1.0)

    if msg is None:
        continue
    if msg.error():
        print(f"⚠️ Kafka error: {msg.error()}")
        continue
    print("🐞 Raw Kafka message:", msg.value().decode("utf-8"))

    try:
        ticket = json.loads(msg.value().decode("utf-8"))
        subject = ticket.get("text", "")
        body = ticket.get("text", "")
        print(f"📩 Received ticket: {subject[:50]}...")
        

        category = predict_category(subject, body)
        print("---------", ticket.get("id"))

        result = {
            "ticket_id": ticket.get("id"),
            "subject": subject,
            "body": body,
            "predicted_category": category,
            "processed_at": time.strftime("%Y-%m-%d %H:%M:%S")
        }
        

        producer.produce(OUTPUT_TOPIC, json.dumps(result).encode("utf-8"))
        producer.flush()
        print(f"✅ Sent classified ticket to {OUTPUT_TOPIC}\n")

    except Exception as e:
        print("❌ Error processing ticket:", e)
