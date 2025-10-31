from confluent_kafka import Consumer
from utils.logger import setup_logger

logger = setup_logger("consumer")

class KafkaConsumer:
    def __init__(self, broker, topic, group_id):
        self.consumer = Consumer({
            "bootstrap.servers": broker,
            "group.id": group_id,
            "auto.offset.reset": "earliest"
        })
        self.consumer.subscribe([topic])

    def poll_message(self, timeout=1.0):
        msg = self.consumer.poll(timeout)
        if msg is None:
            return None
        if msg.error():
            logger.error("Kafka error: %s", msg.error())
            return None
        return msg

    def close(self):
        self.consumer.close()
