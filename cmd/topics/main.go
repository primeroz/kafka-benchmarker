package main

import (
	"fmt"
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
)

var (
	brokerList  = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic       = kingpin.Flag("topic", "Topic name").Default("topic").String()
	nPartitions = kingpin.Flag("nPartitions", "Number Of partitions per topic to create").Default("1").Int32()
	nReplicas   = kingpin.Flag("nReplicas", "Number Of Replicas per topic to create").Default("3").Int16()
	nTopics     = kingpin.Flag("nTopics", "Number Of topics to create").Default("5").Int()
	delete      = kingpin.Flag("delete", "Enable delete mode. default is create").Bool()
)

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin(*brokerList, config)
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = admin.Close() }()

	for i := 0; i < *nTopics; i++ {
		var topicName string

		if *nTopics == 1 {
			topicName = *topic
		} else {
			topicName = fmt.Sprintf("%s-%d", *topic, i)
		}

		if !*delete {
			compression := "snappy"
			minISR := fmt.Sprintf("%d", *nReplicas-1)
			retentionMs := "1800000" // 30 minutes
			segmentBytes := "100000000"
			retentionBytes := "500000000"

			err = admin.CreateTopic(topicName, &sarama.TopicDetail{
				NumPartitions:     *nPartitions,
				ReplicationFactor: *nReplicas,
				ConfigEntries: map[string]*string{
					"compression.type":    &compression,
					"min.insync.replicas": &minISR,
					"retention.ms":        &retentionMs,
					"segment.bytes":       &segmentBytes,
					"retention.bytes":     &retentionBytes,
				},
			}, false)
			if err != nil {
				log.Printf("Error while creating topic: ", err.Error())
			}
		} else {
			err = admin.DeleteTopic(topicName)
			if err != nil {
				log.Printf("Error while creating topic: ", err.Error())
			}
		}
	}
}
