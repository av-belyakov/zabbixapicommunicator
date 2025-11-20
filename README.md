# Zabbixapicommunicator

Пакет '**zabbixapicommunicator**' реализует подключение и дальнейшее взаимодействие с Zabbix с целью передачи данных, получения данных или настройки некоторых действий, которые будут выполнятся Zabbix при изменении данных. Например, действий создания, изменения или удаления групп хостов и хостов, выполнении тригеров, 'действий обнаружения', 'действий авторегистрации' и т.д.

Пакет актуален для версии Zabbix 7.x.

## Состав пакета

Пакет '**zabbixapicommunicator**' состоит из двух подпакетов: '_connectionjsonrpc_' и '_connectionzabbixagent_'.

### Подпакет 'connectionjsonrpc'

Подпакет '_connectionjsonrpc_' позволяет взаимодействовать с API Zabbix по протоколам HTTP или HTTPS методом JsonRPC. Подпакет позволяет получать, обновлять и удалять некоторые данные из базы данных Zabbix. Реализована лишь малая часть монипляций с данными в Zabbix с использованием его API. Подробнее по работе с API Zabbix https://www.zabbix.com/documentation/current/en/manual/api.

Однако, следует обратить внимание, что в подпакете '_connectionjsonrpc_' есть метод **CustomRequest** который позволяет гибко настраивать запросы к API Zabbix, что покрывает большую часть команд API.

### Подпакет 'connectionzabbixagent'

Подпакет '_connectionzabbixagent_' позволяет отправлять некоторые данные Zabbix серверу через официального zabbix-агента.

## Установка пакета

Для того что бы установить пакет версии v2 нужно выполнить:

```bash
go get github.com/av-belyakov/zabbixapicommunicator/v2
```

## Быстрый старт для подпакета 'connectionjsonrpc'

### Подключение

Выполнить подключение к API Zabbix. Пример:

```go
zc, err := connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithPort(80),
			connectionjsonrpc.WithHost("localhost"),
			connectionjsonrpc.WithConnectionTimeout(30),
			connectionjsonrpc.WithLogin("user"),
			connectionjsonrpc.WithPasswd("passwd"),
		)

err = zc.AuthorizationStart(context.Background())
```

### Универсальный метод

Можно использовать универсальный метод, позволяющий выполнят действия над большинством объектов в Zabbix. Однако, этот метод на ряду с гибкостью, требует большей внимательности при составлении запросов. Кроме того под него нет универсального декодера.

```go
bytes, err := zc.CustomRequest(context.Background(), method, request)
```

### Получение данных

Получить список групп хостов и декодировать ответ:

```go
res, err := zc.GetFullHostGroupList(context.Background())
data, errMsg, err := connectionjsonrpc.NewResponseGetHostGroupList().Get(res)
```

Получить список всех хостов для определённой группы и декодировать ответ:

```go
res, err := zc.GetHostList(context.Background(), groupsId...)
data, errMsg, err := connectionjsonrpc.NewResponseGetHostList().Get(res)
```

Получить список всех хостов и декодировать ответ:

```go
res, err := zc.GetHosts(context.Background())
res, errMsg, err := connectionjsonrpc.NewResponseGetHostList().Get(data)
```

### Добавление новых данных

Добавить новую группу хостов и декодировать ответ:

```go
res, err := zc.CreateHostGroup(context.Background(), "имя группы")
data, errMsg, err := connectionjsonrpc.NewResponseCreateHostGroup().Get(res)
```

Добавить новые хосты в группу хостов и декодировать ответ:

```go
res, err := zc.CreateHost(context.Background(), connectionjsonrpc.CreateHostOptionsRequest{})
data, errMsg, err := connectionjsonrpc.NewResponseCreateHost().Get(res)
```

### Обновление данных

Обновить информацию по группе хостов и декодировать ответ:

```go
res, err := zc.GetHostGroup(context.Background(), connectionjsonrpc.FilterHostGroup{})
data, errMsg, err := connectionjsonrpc.NewResponseGetHostGroupList().Get(res)
```

Обновление в хосте параметра 'группы хостов':

```go
res, err := zc.CreateHostGroup(context.Background(), hostGroupName)
```

Обновление в хосте параметра 'теги':

```go
res, err := zc.UpdateHostParameterTags(context.Background(), hostId, connectionjsonrpc.Tags{})
```

Обновление в хосте параметра 'макросы':

```go
res, err := zc.UpdateHostParameterMacro(context.Background(), hostId, connectionjsonrpc.Macros{})
```

Обновление в хосте параметра 'интерфейсы'

```go
res, err := zc.UpdateHostParameterInterfaces(context.Background(), hostId, connectionjsonrpc.InterfacesRequest{})
```

Запись нового параметра 'инвентаризация', старые данные будут затёрты:

```go
res, err := zc.CreateHostParameterInventory(context.Background(), hostId, connectionjsonrpc.HostInventory{})
```

### Удаление данных

Удаление хоста:

```go
res, err := zc.DeleteHost(context.Background(), hostId...)
```

Удаление группы хостов:

```go
res, err = zc.DeleteHostGroup(context.Background(), groupsId...)
```
