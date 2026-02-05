package mqtt

import (
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	pahoClient paho.Client
}

type Config struct {
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
}

func NewClient(cfg Config) *Client {
	opts := paho.NewClientOptions()
	opts.AddBroker(cfg.BrokerURL)
	opts.SetClientID(cfg.ClientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetOnConnectHandler(func(c paho.Client) {
		log.Printf("MQTT: Connected to broker at %s", cfg.BrokerURL)
	})
	opts.SetConnectionLostHandler(func(c paho.Client, err error) {
		log.Printf("MQTT: Connection lost: %v", err)
	})

	return &Client{
		pahoClient: paho.NewClient(opts),
	}
}

func (c *Client) Connect() error {
	token := c.pahoClient.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connect: %w", token.Error())
	}
	return nil
}

func (c *Client) Disconnect() {
	c.pahoClient.Disconnect(250)
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.pahoClient.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT publish: %w", token.Error())
	}
	return nil
}
