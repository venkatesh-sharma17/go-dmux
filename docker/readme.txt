./kafka-topics.sh --create --zookeeper zookeeper:2181/kafka --replication-factor 1 --partitions 1 --topic sample-topic

./kafka-topics.sh --list --zookeeper zookeeper:2181/kafka

./kafka-topics.sh --describe --topic sample-topic --zookeeper zookeeper:2181/kafka

./kafka-console-producer.sh --broker-list <host>:9092 --topic sample-topic
    
echo "key1:value1"| /opt/kafka_2.11-2.0.1/bin/kafka-console-producer.sh --broker-list <host>:9092 --topic sample-topic --property "parse.key=true" --property "key.separator=:"

./kafka-console-consumer.sh --zookeeper zookeeper:2181/kafka --topic test --from-beginning

/opt/kafka_2.11-2.0.1/bin/kafka-console-consumer.sh --bootstrap-server <host>:9092 --topic sample-topic --from-beginning