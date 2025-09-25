package mqttchat

import (
	"encoding/json"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mateusxxlima/chat-mqtt/internal/config"
)

const QOS = byte(2)

func NewMQTTClient() (conn MQTT.Client) {
	MOSQUITTO_HOST := config.Env.MosquittoHost
	opts := MQTT.NewClientOptions().AddBroker(MOSQUITTO_HOST)
	opts.SetClientID(SELF_ID)
	opts.SetDefaultPublishHandler(sysMsgListener)
	opts.CleanSession = true
	conn = MQTT.NewClient(opts)
	if token := conn.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return
}

func pubInTopic(conn MQTT.Client, topic string, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error on parsing data to json:", err)
		fmt.Println("Error on publishing msg to topic", topic)
		return
	}
	token := conn.Publish(topic, QOS, false, payload)
	token.Wait()
}

func subInTopic(conn MQTT.Client, topic string, callback MQTT.MessageHandler) {
	if token := conn.Subscribe(topic, QOS, callback); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
