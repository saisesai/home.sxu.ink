package mqtt

import (
	"errors"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/saisesai/home.sxu.ink/config"
	"github.com/saisesai/home.sxu.ink/log"
)

var Client paho.Client

var ErrInvalidCmd = errors.New("invalid cmd")

var ErrDeviceNotExist = errors.New("device not exist")

func init() {
	opts := paho.NewClientOptions()
	opts.KeepAlive = 60
	opts.AddBroker(config.MqttBrokerUrl)
	opts.SetClientID(config.ServerMqttClientId)
	opts.SetUsername(config.MqttUsername)
	opts.SetPassword(config.MqttPassword)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = paho.NewClient(opts)
}

var connectHandler paho.OnConnectHandler = func(client paho.Client) {
	L := log.WithField("from", "mqtt")
	L.Infoln("broker connected!")
	err := subscribeDevices()
	if err != nil {
		log.WithField("error", err).Fatalln("failed to subscribe devices!")
	}
}

var connectLostHandler paho.ConnectionLostHandler = func(client paho.Client, err error) {
	L := log.WithField("from", "mqtt")
	L.Warningln("lost connection with broker!")
}

func subscribeDevices() error {
	var err error
	// relay
	err = subscribeRelays()
	if err != nil {
		return err
	}
	// ...
	return nil
}
