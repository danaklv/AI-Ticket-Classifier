from confluent_kafka import Producer
import json
from utils.logger import setup_logger

logger = setup_logger("producer")

class KafkaProducer:
    def __init__(self, broker):
        self.producer = Producer({
            "bootstrap.servers": broker,
            "linger.ms": 50,      
            "batch.num.messages": 10000,  
            "queue.buffering.max.messages": 100000,  
            "enable.idempotence": True     
        })

    def delivery_callback(self, err, msg):
        if err is not None:
            logger.error(f"Delivery failed for record {msg.key()}: {err}")
        else:
            logger.info(
                f"Message delivered to {msg.topic()} [{msg.partition()}] @ offset {msg.offset()}"
            )

    def send(self, topic, message):
        try:
            self.producer.produce(
                topic,
                json.dumps(message).encode("utf-8"),
                callback=self.delivery_callback
            )
            self.producer.poll(0)  
        except BufferError:
            logger.warning("Producer queue is full. Flushing...")
            self.producer.flush(5)
        except Exception as e:
            logger.error(f"Failed to send message: {e}")

    def flush(self):
        logger.info("Flushing producer queue...")
        self.producer.flush()