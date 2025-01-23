package zabbixapicommunicator_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	zabbixapicommunicator "github.com/av-belyakov/zabbixapicommunicator/cmd"
)

type Message struct {
	Type, Message string
}

func (m *Message) GetType() string {
	return m.Type
}

func (m *Message) GetMessage() string {
	return m.Message
}

func (m *Message) SetType(v string) {
	m.Type = v
}

func (m *Message) SetMessage(v string) {
	m.Message = v
}

func TestMethodStart(t *testing.T) {
	z, err := zabbixapicommunicator.New(zabbixapicommunicator.SettingsZabbixConnection{
		Host:       "127.0.0.1",
		Port:       55441,
		ZabbixHost: "app.test.z",
	})
	assert.NoError(t, err)

	chRecipient := make(chan zabbixapicommunicator.Messager)
	ctx, ctxDone := context.WithCancel(context.Background())
	z.Start(ctx, []zabbixapicommunicator.EventType{
		{
			IsTransmit: true,
			EventType:  "error",
			ZabbixKey:  "f83rhrurr-rr484h",
		},
		{
			IsTransmit: true,
			EventType:  "info",
			ZabbixKey:  "7373r8fef-ff3uf4",
		},
	}, chRecipient)

	newMessageSettings := &zabbixapicommunicator.MessageSettings{}
	newMessageSettings.SetType("error")
	newMessageSettings.SetMessage("some error message")
	chRecipient <- newMessageSettings

	go func() {
		time.Sleep(3 * time.Second)
		ctxDone()
	}()

	for err := range z.GetChanErr() {
		t.Log("error", err)
	}

	assert.True(t, true)
}
