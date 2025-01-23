package zabbixapicommunicator

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Start обработчик запросов, от внешних модулей, которые необходимо передать в Zabbix
func (zc *ZabbixConnection) Start(ctx context.Context, events []EventType, chRecipient <-chan Messager) error {
	if len(events) == 0 {
		return errors.New("invalid configuration file for Zabbix, the number of event types (ZABBIX.zabbixHosts.eventTypes) is 0")
	}

	go func() {
		var wg sync.WaitGroup
		listChans := make(map[string]chan<- string, len(events))

		for _, v := range events {
			if !v.IsTransmit {
				continue
			}

			newChan := make(chan string)
			listChans[v.EventType] = newChan

			wg.Add(1)
			go func(chMsg <-chan string, zKey string, handshake Handshake) {
				defer wg.Done()

				var ticker *time.Ticker
				if handshake.TimeInterval > 0 && handshake.Message != "" {
					ticker = time.NewTicker(time.Duration(handshake.TimeInterval) * time.Minute)
					defer ticker.Stop()
				}

				if ticker == nil {
					for msg := range chMsg {
						if _, err := zc.sendData(zKey, []string{msg}); err != nil {
							zc.chanErr <- err
						}
					}
				} else {
					for {
						select {
						case <-ticker.C:
							if _, err := zc.sendData(zKey, []string{handshake.Message}); err != nil {
								zc.chanErr <- err
							}

						case msg, open := <-chMsg:
							if !open {
								return
							}

							if _, err := zc.sendData(zKey, []string{msg}); err != nil {
								zc.chanErr <- err
							}
						}
					}
				}
			}(newChan, v.ZabbixKey, v.Handshake)
		}

		go func() {
			defer func() {
				for _, channel := range listChans {
					close(channel)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return

				case msg, open := <-chRecipient:
					if !open {
						return
					}

					if ch, ok := listChans[msg.GetType()]; ok {
						ch <- msg.GetMessage()
					}
				}
			}
		}()

		wg.Wait()

		close(zc.chanErr)
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
func (zc *ZabbixConnection) sendData(zkey string, data []string) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("the list of transmitted data should not be empty")
	}

	ldz := make([]DataZabbix, 0, len(data))
	for _, v := range data {
		ldz = append(ldz, DataZabbix{
			Host:  zc.zabbixHost,
			Key:   zkey,
			Value: v,
		})
	}

	jsonReg, err := json.Marshal(PatternZabbix{
		Request: "sender data",
		Data:    ldz,
	})
	if err != nil {
		return 0, err
	}

	//заголовок пакета
	pkg := []byte("ZBXD\x01")

	//длинна пакета с данными
	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(jsonReg)))

	pkg = append(pkg, dataLen...)
	pkg = append(pkg, jsonReg...)

	var d net.Dialer = net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), *zc.connTimeout)
	defer cancel()

	conn, err := d.DialContext(ctx, zc.netProto, fmt.Sprintf("%s:%d", zc.host, zc.port))
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	num, err := conn.Write(pkg)
	if err != nil {
		return 0, err
	}

	_, err = io.ReadAll(conn)
	if err != nil {
		return num, err
	}

	return num, nil
}
