package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
)

var (
	brokerList      = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	topic           = kingpin.Flag("topic", "Topic name").Default("topic").String()
	topicRangeStart = kingpin.Flag("topicRangeStart", "Topic range Start").Default("0").Int()
	topicRangeEnd   = kingpin.Flag("topicRangeEnd", "Topic range End").Default("0").Int()
	nMessages       = kingpin.Flag("nMessages", "Number of Messages").Default("1000").Int()
	nThreads        = kingpin.Flag("nThreads", "Number of Threads").Default("3").Int()
	maxRetry        = kingpin.Flag("maxRetry", "Retry limit").Default("1").Int()
)

var wg sync.WaitGroup

func RandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func produceInRandomTopic(producer sarama.SyncProducer, messages int) {
	defer wg.Done()

	var topicName string

	for m := 0; m < messages; m++ {
		if *topicRangeStart == 0 && *topicRangeEnd == 0 {
			topicName = *topic
		} else if *topicRangeStart == *topicRangeEnd {
			topicName = fmt.Sprintf("%s-%d", *topic, *topicRangeStart)
		} else {
			topicName = fmt.Sprintf("%s-%d", *topic, rand.Intn(*topicRangeEnd-*topicRangeStart)+(*topicRangeStart))
			//topicName = fmt.Sprintf("%s-%d", *topic, int(math.Min((rand.ExpFloat64()/10)*max, max)))
		}

		//myStr := RandomString(int(math.Min((math.Abs(rand.NormFloat64()*float64(1000000/2)) + float64(1000000/2)), 999999)))
		strLength := math.Min((rand.ExpFloat64()/10)*1000000, 999999)
		myStr := RandomString(int(strLength))

		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Value: sarama.StringEncoder(myStr),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("%s", err)
		}
		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d) of length %d\n", topicName, partition, offset, len(myStr))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *maxRetry
	config.Producer.Return.Successes = true
	//config.Producer.Retry.Backoff = (5 * time.Second)
	//config.Producer.Timeout = (5 * time.Second)
	//config.Net.ReadTimeout = (5 * time.Second)
	//config.Net.DialTimeout = (5 * time.Second)
	//config.Net.WriteTimeout = (5 * time.Second)
	//config.Metadata.Retry.Max = 1
	//config.Metadata.Retry.Backoff = (1 * time.Second)
	//config.Metadata.RefreshFrequency = (5 * time.Second)
	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	runtime.GOMAXPROCS(*nThreads)
	for i := 0; i < *nThreads; i++ {
		wg.Add(1)
		go produceInRandomTopic(producer, (*nMessages)/(*nThreads))
	}
	wg.Wait()
}
