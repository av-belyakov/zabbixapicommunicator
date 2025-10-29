# Zabbixapicommunicator

Пакет '**zabbixapicommunicator**' реализует подключение и дальнейшее взаимодействие с Zabbix с целью передачи данных, получения данных или настройки некоторых действий которые будут выполнятся Zabbix при изменении данных. Например, таких действий как выполнение тригеров, 'действий обнаружения', 'действий авторегистрации' и т.д.

Пакет актуален для версии Zabbix 7.x.

Пакет '**zabbixapicommunicator**' состоит из двух подпакетов: _connectionjsonrpc_ и _connectionzabbixagent_.

### Пакет _connectionjsonrpc_

Пакет _connectionjsonrpc_ позволяет взаимодействовать с Zabbix сервером по протоколам HTTP или HTTPS путем обращения к модулю Zabbix реализующий взаимодействие матодом JsonRPC. С помощью этого пакета можно запрашивать данные из базы данных Zabbix, добавлять, удалять триггеры, 'действия обнаружения', 'действия авторегистрации' и т.д. Подробнее в документайии https://www.zabbix.com/documentation/current/en/manual/api.

Подробнее по работе с API Zabbix https://www.zabbix.com/documentation/current/en/manual/api.

### Пакет _connectionzabbixagent_

Пакет _connectionzabbixagent_ позволяет отправлять некоторые данные Zabbix серверу через его официального zabbix-агента.
