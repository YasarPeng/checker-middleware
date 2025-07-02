package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"checker-middleware/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

func ensureKafkaTopic(cfg MQConfig) error {
	// 只取第一个broker用于创建topic
	broker := cfg.Brokers[0]
	conn, err := kafka.Dial("tcp", broker)
	logger.DebugLog("ensureKafkaTopic: Brokers[0] %s", broker)
	if err != nil {
		return fmt.Errorf("kafka dial error: %v", err)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	logger.DebugLog("Kafka controller: host=%s port=%d", controller.Host, controller.Port)
	if err != nil {
		return fmt.Errorf("kafka controller error: %v", err)
	}
	var controllerConn *kafka.Conn
	controllerAddr := net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port))
	controllerConn, err = kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("kafka controller dial error: %v", err)
	}
	defer controllerConn.Close()
	topicConfigs := []kafka.TopicConfig{{
		Topic:             cfg.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}}
	return controllerConn.CreateTopics(topicConfigs...)
}

func MQConnect(cfg MQConfig) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "kafka":
		// 自动创建topic
		if err := ensureKafkaTopic(cfg); err != nil {
			result["error"] = fmt.Sprintf("kafka create topic error: %v", err)
			return result
		}
		r := kafka.ReaderConfig{
			Brokers:   cfg.Brokers,
			Topic:     cfg.Topic,
			Partition: 0,
			MinBytes:  1,
			MaxBytes:  10e6,
		}
		reader := kafka.NewReader(r)
		defer reader.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := reader.ReadLag(ctx)
		if err != nil {
			result["error"] = fmt.Sprintf("kafka connect error: %v", err)
			return result
		}
		result["success"] = "true"
	case "rabbitmq", "mq":
		// 检查vhost是否存在，连接时带vhost
		vhost := cfg.Vhost
		if vhost == "" {
			vhost = "/" // 默认vhost
		}
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, vhost)
		logger.DebugLog("rabbitmq connect url: %s", url)
		conn, err := amqp.Dial(url)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq connect error (vhost=%s): %v", vhost, err)
			return result
		}
		defer conn.Close()
		ch, err := conn.Channel()
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq channel error: %v", err)
			return result
		}
		defer ch.Close()
		_, err = ch.QueueDeclare("laiye_precheck", true, false, false, false, nil)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq queue error: %v", err)
			return result
		}
		result["success"] = "true"
	default:
		result["error"] = "unsupported provider"
	}

	return result
}

func MQWrite(cfg MQConfig, msg string) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "kafka":
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
		})
		defer writer.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := writer.WriteMessages(ctx, kafka.Message{Value: []byte(msg)})
		if err != nil {
			result["error"] = fmt.Sprintf("kafka write error: %v", err)
			return result
		}
		result["success"] = "true"
	case "rabbitmq":
		vhost := cfg.Vhost
		if vhost == "" {
			vhost = "/"
		}
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, vhost)
		conn, err := amqp.Dial(url)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq connect error (vhost=%s): %v", vhost, err)
			return result
		}
		defer conn.Close()
		ch, err := conn.Channel()
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq channel error: %v", err)
			return result
		}
		defer ch.Close()
		err = ch.Publish("", "laiye_precheck", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
		logger.DebugLog("rabbitmq write msg: %s", msg)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq write error: %v", err)
			return result
		}
		result["success"] = "true"
	default:
		result["error"] = "unsupported provider"
	}
	return result
}

func MQDelete(cfg MQConfig) map[string]string {
	result := map[string]string{"success": "false"}
	provider := strings.ToLower(cfg.Provider)
	switch provider {
	case "kafka":
		// Kafka 没有直接删除消息的API，这里用消费一条消息模拟“删除”
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: "precheck-delete",
		})
		defer reader.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		msg, err := reader.ReadMessage(ctx)
		logger.DebugLog("Kafka Consumed message: topic=%s partition=%d offset=%d key=%s value=%s", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		if err != nil {
			result["error"] = fmt.Sprintf("kafka consume(delete) error: %v", err)
			return result
		}
		result["success"] = "true"
	case "rabbitmq":
		vhost := cfg.Vhost
		if vhost == "" {
			vhost = "/"
		}
		queueName := "laiye_precheck"
		if queueName == "" {
			queueName = "default"
		}
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, vhost)
		conn, err := amqp.Dial(url)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq connect error: %v", err)
			return result
		}
		defer conn.Close()
		ch, err := conn.Channel()
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq channel error: %v", err)
			return result
		}
		defer ch.Close()
		// 先消费一条消息（如果有）
		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err == nil {
			select {
			case d := <-msgs:
				_ = d.Ack(false)
				// 消费成功，不管ack是否报错都继续
				logger.DebugLog("rabbitmq consume success")
			case <-time.After(5 * time.Second):
				logger.DebugLog("rabbitmq no message to delete")
				// 没有消息，继续
			}
		}
		// 再删除队列
		_, err = ch.QueueDelete(
			queueName,
			false,
			false,
			false,
		)
		if err != nil {
			result["error"] = fmt.Sprintf("rabbitmq queue delete error: %v", err)
			return result
		}
		result["success"] = "true"
	default:
		result["error"] = "unsupported provider"
	}
	return result
}

func VerifyMQ(cfg MQConfig) MQResult {
	content := "hello"
	res := MQResult{
		Connect: MQConnect(cfg),
		Write:   map[string]string{"success": "skip"},
		Delete:  map[string]string{"success": "skip"},
	}
	res.Write = MQWrite(cfg, content)
	if res.Connect["success"] == "true" {
		res.Write = MQWrite(cfg, content)
		if res.Write["success"] == "true" {
			res.Delete = MQDelete(cfg)
		}
	}
	return res
}

func VerifyMQJson(cfg MQConfig) []byte {
	res := VerifyMQ(cfg)
	b, _ := json.Marshal(res)
	return b
}
