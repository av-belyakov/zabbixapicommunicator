package zabbixapicommunicator

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

// Start запуск обработчика заапросов к Zabbix
func (zc *ZabbixConnection) Start(ctx context.Context, events []EventType, recipient <-chan Messager) error {
	countEvents := len(events)
	if countEvents == 0 {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'invalid configuration file for Zabbix, the number of event types (ZABBIX.zabbixHosts.eventTypes) is 0' %s:%d", f, l-1)
	}

	listChans := make(map[string]chan<- string, countEvents)

	go func() {
		<-ctx.Done()

		for _, channel := range listChans {
			close(channel)
		}
		listChans = nil

		close(zc.chanErr)
	}()

	for _, v := range events {
		if !v.IsTransmit {
			continue
		}

		newChan := make(chan string)
		listChans[v.EventType] = newChan

		go func(cm <-chan string, zkey string, hs Handshake) {
			var t *time.Ticker
			if hs.TimeInterval > 0 && hs.Message != "" {
				t = time.NewTicker(time.Duration(hs.TimeInterval) * time.Minute)
				defer t.Stop()
			}

			if t == nil {
				for msg := range cm {
					if _, err := zc.sendData(zkey, []string{msg}); err != nil {
						zc.chanErr <- err
					}
				}
			} else {
				for {
					select {
					case <-t.C:
						if _, err := zc.sendData(zkey, []string{hs.Message}); err != nil {
							zc.chanErr <- err
						}

					case msg, open := <-cm:
						if !open {
							cm = nil

							return
						}

						if _, err := zc.sendData(zkey, []string{msg}); err != nil {
							zc.chanErr <- err
						}
					}
				}
			}
		}(newChan, v.ZabbixKey, v.Handshake)
	}

	go func() {
		for msg := range recipient {
			if c, ok := listChans[msg.GetType()]; ok {
				c <- msg.GetMessage()
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
