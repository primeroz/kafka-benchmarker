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
	nPartitions = kingpin.Flag("nPartitions", "Number Of partitions per topic to create").Default("1").Int()
	nReplicas   = kingpin.Flag("nReplicas", "Number Of Replicas per topic to create").Default("3").Int()
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
			topicName = fmt.Sprintf("%s-%s", *topic, i)
		}

		if !*delete {
			err = admin.CreateTopic(topicName, &sarama.TopicDetail{
				NumPartitions:     *nPartitions,
				ReplicationFactor: *nReplicas,
			}, false)
			if err != nil {
				log.Printf("Error while creating topic: ", err.Error())
			}
		} else {
			err = admin.CreateTopic(topicName)
			if err != nil {
				log.Printf("Error while creating topic: ", err.Error())
			}
		}
	}
}
