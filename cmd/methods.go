package zabbixapicommunicator

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

// Start обработчик запросов, от внешних модулей, которые необходимо передать в Zabbix
func (zc *ZabbixConnection) Start(ctx context.Context, events []EventType, chRecipient <-chan Messager) error {
	if len(events) == 0 {
		return errors.New("invalid configuration file for Zabbix, the number of event types (ZABBIX.zabbixHosts.eventTypes) is 0")
	}

	listChans := map[string]chan string{}

	for _, v := range events {
		if !v.IsTransmit {
			continue
		}

		newChan := make(chan string)
		listChans[v.EventType] = newChan

		go func(ctx context.Context, key string, handshake Handshake, ch <-chan string) {
			var ticker *time.Ticker
			if handshake.TimeInterval > 0 {
				ticker = time.NewTicker(time.Duration(handshake.TimeInterval) * time.Minute)
				defer ticker.Stop()
			}

			if ticker == nil {
				for msg := range ch {
					go zc.sendData(ctx, key, []string{msg})
				}

				return
			}

			for {
				select {
				case <-ctx.Done():
					return

				case <-ticker.C:
					go zc.sendData(ctx, key, []string{handshake.Message})

				case msg, open := <-ch:
					if !open {
						return
					}

					go zc.sendData(ctx, key, []string{msg})
				}
			}
		}(ctx, v.ZabbixKey, v.Handshake, newChan)
	}

	//обработка входящих сообщений
	go func() {
		defer func() {
			for _, channel := range listChans {
				if channel == nil {
					continue
				}

				close(channel)
			}

			close(zc.chanErr)
		}()

		for {
			select {
			case <-ctx.Done():
				zc.chanErr <- errors.New("context stop")

				return

			case msg, open := <-chRecipient:
				if !open {
					zc.chanErr <- errors.New("the channel for receiving messages coming to the module has been closed")

					return
				}

				if ch, ok := listChans[msg.GetType()]; ok {
					if ch == nil {
						continue
					}

					ch <- msg.GetMessage()
				}
			}
		}
	}()

	return nil
}

// GetChanErr метод возвращающий канал в который отправляются ошибки возникающие при соединении с Zabbix
func (zc *ZabbixConnection) GetChanErr() chan error {
	return zc.chanErr
}

// GetType возвращает тип события
func (s *MessageSettings) GetType() string {
	return s.EventType
}

// SetType устанавливает тип события
func (s *MessageSettings) SetType(v string) {
	s.EventType = v
}

// GetMessage возвращает сообщение
func (s *MessageSettings) GetMessage() string {
	return s.Message
}

// SetMessage устанавливает сообщение
func (s *MessageSettings) SetMessage(v string) {
	s.Message = v
}

// sendData метод реализующий отправку данных в Zabbix
func (zc *ZabbixConnection) sendData(ctx context.Context, key string, data []string) {
	if len(data) == 0 {
		zc.chanErr <- errors.New("the list of transmitted data should not be empty")

		return
	}

	dataZabbix := make([]DataZabbix, len(data))
	for _, v := range data {
		dataZabbix = append(dataZabbix, DataZabbix{
			Host:  zc.zabbixHost,
			Key:   key,
			Value: v,
		})
	}

	reg, err := json.Marshal(PatternZabbix{
		Request: "sender data",
		Data:    dataZabbix,
	})
	if err != nil {
		zc.chanErr <- err

		return
	}

	//заголовок пакета
	pkg := []byte("ZBXD\x01")

	//длинна пакета с данными
	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(reg)))

	pkg = append(pkg, dataLen...)
	pkg = append(pkg, reg...)

	var d net.Dialer = net.Dialer{Timeout: zc.connTimeout}
	conn, err := d.DialContext(ctx, zc.netProto, fmt.Sprintf("%s:%d", zc.host, zc.port))
	if err != nil {
		zc.chanErr <- err

		return
	}
	defer conn.Close()

	_, err = conn.Write(pkg)
	if err != nil {
		zc.chanErr <- err

		return
	}
}
