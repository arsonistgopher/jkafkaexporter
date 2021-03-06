package kafka

import (
	"strconv"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/arsonistgopher/jkafkaexporter/internal/channels"
	"github.com/arsonistgopher/jkafkaexporter/junoscollector"
)

// Config holds the Kafka info
type Config struct {
	KafkaExport time.Duration // Number of seconds in-between kafka exports
	KafkaHost   string        // Host IP or FQDN of kafka bus
}

// StartKafka is a GR that accepts a channel.
func StartKafka(me string, kc Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup) (chan channels.Response, error) {

	// Create the buffer depth to match the number of collectors
	responsechan := make(chan channels.Response, jc.Len())

	go func(me string, kc Config, jc *junoscollector.JunosCollector, done chan struct{}, wg *sync.WaitGroup, responsechan chan channels.Response) {
		ticker := time.NewTicker(kc.KafkaExport)
		kafkadeath := make(chan struct{}, 1)

		var wg2 sync.WaitGroup
		// Add one on for the Go routine we're about to launch
		wg2.Add(1)

		go func(responsechan chan channels.Response, kc Config, done chan struct{}, wg2 *sync.WaitGroup) {

			config := sarama.NewConfig()
			config.Producer.Retry.Max = 5
			config.Producer.RequiredAcks = sarama.WaitForAll
			config.Producer.Return.Successes = true
			brokers := []string{kc.KafkaHost}

			pd, err := sarama.NewSyncProducer(brokers, config)
			if err != nil {
				// Should not reach here
				panic(err)
			}

			defer func() {
				if err := pd.Close(); err != nil {
					// Should not reach here
					panic(err)
				}
			}()
			for {
				select {
				case <-done:
					wg2.Done()
					return
				case r := <-responsechan:
					strTime := strconv.Itoa(int(time.Now().Unix()))
					msg := &sarama.ProducerMessage{
						Topic: r.Topic,
						Key:   sarama.StringEncoder(strTime),
						Value: sarama.StringEncoder(r.Data),
					}

					_, _, err = pd.SendMessage(msg)
					if err != nil {
						panic(err)
					}
				}
			}
		}(responsechan, kc, kafkadeath, &wg2)

		for {
			select {
			case <-done:
				// Get's here, we're done
				kafkadeath <- struct{}{}
				wg2.Wait()
				wg.Done()
				return
			case <-ticker.C:
				// For each collector item, collect and dump
				jc.Collect(responsechan, me)
			}
		}
	}(me, kc, jc, done, wg, responsechan)

	return responsechan, nil
}
