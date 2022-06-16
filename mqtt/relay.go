package mqtt

import (
	"encoding/json"
	"errors"
	paho "github.com/eclipse/paho.mqtt.golang"
	log "github.com/saisesai/home.sxu.ink/log"
	"github.com/saisesai/home.sxu.ink/model"
	"gorm.io/gorm"
)

func subscribeRelays() error {
	var err error
	L := log.WithField("from", "mqtt")
	clientIds, err := model.GetAllRelayClientId()
	if err != nil { // internal error
		return err
	}
	L.Infof("get %d relays from db.", len(clientIds))
	for _, v := range clientIds {
		token := Client.Subscribe(v+"/msg", 0, RelayMsgHandler)
		token.Wait()
		if token.Error() != nil {
			L.Errorln("failed to subscribe topic: ", v+"/msg")
			return err
		}
		token = Client.Subscribe(v+"/oln", 0, RelayOnlineHandler)
		token.Wait()
		if token.Error() != nil {
			L.Errorln("failed to subscribe topic: ", v+"/oln")
			return err
		}
	}
	return nil
}

func RelayOnlineHandler(pClient paho.Client, pMsg paho.Message) {
	var err error
	topic := pMsg.Topic()     // xxx/msg
	payload := pMsg.Payload() // json
	clientId := topic[:len(topic)-4]
	L := log.
		WithField("device", "relay").
		WithField("topic", topic).
		WithField("payload", string(payload))
	L.Debugln("relay msg received!")
	// parse payload
	var pm map[string]interface{}
	err = json.Unmarshal(payload, &pm)
	if err != nil {
		L.Warningln("invalid data format!")
		return
	}
	// find relay
	relay, err := model.FindRelayByClientId(clientId)
	// check field online
	if online, ok := pm["online"].(bool); ok {
		relay.Online = online
	}
	err = relay.Save()
	if err != nil {
		L.Errorln("failed to save relay info!")
	}
}

func RelayMsgHandler(pClient paho.Client, pMsg paho.Message) {
	var err error
	topic := pMsg.Topic()     // xxx/msg
	payload := pMsg.Payload() // json
	clientId := topic[:len(topic)-4]
	L := log.
		WithField("device", "relay").
		WithField("topic", topic).
		WithField("payload", string(payload))
	L.Debugln("relay msg received!")
	// parse payload
	var pm map[string]interface{}
	err = json.Unmarshal(payload, &pm)
	if err != nil {
		L.Warningln("invalid data format!")
		return
	}
	// find relay
	relay, err := model.FindRelayByClientId(clientId)
	// check field hv, sv, online, on
	if hv, ok := pm["hv"].(string); ok {
		relay.HardwareVersion = hv
	}
	if sv, ok := pm["sv"].(string); ok {
		relay.SoftwareVersion = sv
	}
	if on, ok := pm["on"].(bool); ok {
		relay.On = on
	}
	err = relay.Save()
	if err != nil {
		L.Errorln("failed to save relay info!")
	}
}

// SendRelayCmd
// return: bad request(ErrInvalidCmd|ErrDeviceNotExist) or internal error
func SendRelayCmd(pClientId string, pCmd string) error {
	var err error
	// check cmd(ON OFF)
	if pCmd != "ON" && pCmd != "OFF" {
		return ErrInvalidCmd // bad request
	}
	// check relay exist
	_, err = model.FindRelayByClientId(pClientId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDeviceNotExist // bad request
		}
		return err // internal error
	}
	// send cmd
	topic := pClientId + "/cmd"
	if token := Client.Publish(topic, 0, false, pCmd); token.Wait() && token.Error() != nil {
		return token.Error() // internal error
	}
	return nil
}
