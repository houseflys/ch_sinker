/*Copyright [2019] housepower

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package input

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/sundy-li/go_commons/log"

	"github.com/housepower/clickhouse_sinker/config"
	"github.com/housepower/clickhouse_sinker/model"
	"github.com/housepower/clickhouse_sinker/statistics"
)

// KafkaGo implements input.Inputer
type KafkaGo struct {
	taskCfg *config.TaskConfig
	r       *kafka.Reader
	stopped chan struct{}
	putFn   func(msg model.InputMessage)
}

// NewKafkaGo get instance of kafka reader
func NewKafkaGo() *KafkaGo {
	return &KafkaGo{}
}

// Init Initialise the kafka instance with configuration
func (k *KafkaGo) Init(taskCfg *config.TaskConfig, putFn func(msg model.InputMessage)) error {
	k.taskCfg = taskCfg
	cfg := config.GetConfig()
	kfkCfg := cfg.Kafka[k.taskCfg.Kafka]
	k.stopped = make(chan struct{})
	k.putFn = putFn
	if kfkCfg.Sasl.Enable && kfkCfg.Sasl.Username == "" {
		return errors.Errorf("kafka-go doesn't support SASL/GSSAPI(Kerberos)")
	}
	offset := kafka.LastOffset
	if k.taskCfg.Earliest {
		offset = kafka.FirstOffset
	}
	k.r = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(kfkCfg.Brokers, ","),
		GroupID:        k.taskCfg.ConsumerGroup,
		Topic:          k.taskCfg.Topic,
		StartOffset:    offset,
		MinBytes:       k.taskCfg.MinBufferSize * k.taskCfg.MsgSizeHint,
		MaxBytes:       k.taskCfg.BufferSize * k.taskCfg.MsgSizeHint,
		MaxWait:        time.Duration(k.taskCfg.FlushInterval) * time.Second,
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})
	return nil
}

// kafka main loop
func (k *KafkaGo) Run(ctx context.Context) {
LOOP_KAFKA_GO:
	for {
		var err error
		var msg kafka.Message
		if msg, err = k.r.FetchMessage(ctx); err != nil {
			switch errors.Cause(err) {
			case context.Canceled:
				log.Infof("%s: Kafka.Run quit due to context has been canceled", k.taskCfg.Name)
				break LOOP_KAFKA_GO
			case io.EOF:
				log.Infof("%s: Kafka.Run quit due to reader has been closed", k.taskCfg.Name)
				break LOOP_KAFKA_GO
			default:
				statistics.ConsumeMsgsErrorTotal.WithLabelValues(k.taskCfg.Name).Inc()
				err = errors.Wrap(err, "")
				log.Errorf("%s: Kafka.Run got error %+v", k.taskCfg.Name, err)
				continue
			}
		}
		k.putFn(model.InputMessage{
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Key:       msg.Key,
			Value:     msg.Value,
			Offset:    msg.Offset,
			Timestamp: &msg.Time,
		})
	}
}

func (k *KafkaGo) CommitMessages(ctx context.Context, msg *model.InputMessage) error {
	err := k.r.CommitMessages(ctx, kafka.Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
	if err != nil {
		err = errors.Wrapf(err, "")
		return err
	}
	return nil
}

// Stop kafka consumer and close all connections
func (k *KafkaGo) Stop() error {
	if k.r != nil {
		k.r.Close()
	}
	return nil
}

// Description of this kafka consumer, which topic it reads from
func (k *KafkaGo) Description() string {
	return "kafka consumer of topic " + k.taskCfg.Topic
}
