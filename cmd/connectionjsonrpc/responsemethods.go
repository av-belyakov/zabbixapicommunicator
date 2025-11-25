package connectionjsonrpc

import (
	"encoding/json"

	"github.com/av-belyakov/zabbixapicommunicator/v2/internal/supportingfunctions"
)

// NewResponseAPIInfo ответ на запрос с целью получения информации о API
func NewResponseAPIInfo() *ResponseAPIInfo {
	return &ResponseAPIInfo{}
}

// Get получить информацию о API
func (r *ResponseAPIInfo) Get(b []byte) (*ResponseAPIInfo, *ResponseError, error) {
	res := &ResponseAPIInfo{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if res.Result == "" {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseCreateHostGroup ответ на создание группы хостов
func NewResponseCreateHostGroup() *ResponseCreateHostGroup {
	return &ResponseCreateHostGroup{}
}

// Get получить ответ на создание группы хостов
func (r *ResponseCreateHostGroup) Get(b []byte) (*ResponseCreateHostGroup, *ResponseError, error) {
	res := &ResponseCreateHostGroup{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result.GroupIds) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseCreateHost ответ на создание хоста
func NewResponseCreateHost() *ResponseCreateHost {
	return &ResponseCreateHost{}
}

// Get получить ответ на создание хоста
func (r *ResponseCreateHost) Get(b []byte) (*ResponseCreateHost, *ResponseError, error) {
	res := &ResponseCreateHost{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result.HostIds) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseGetHostGroupList ответ на запрос с целью получения списка групп хостов
func NewResponseGetHostGroupList() *ResponseHostGroupList {
	return &ResponseHostGroupList{}
}

// Get получить список групп хостов
func (r *ResponseHostGroupList) Get(b []byte) (*ResponseHostGroupList, *ResponseError, error) {
	res := &ResponseHostGroupList{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseGetHostList ответ на запрос с целью получения списка хостов
func NewResponseGetHostList() *ResponseHostList {
	return &ResponseHostList{}
}

// Get получить список хостов
func (r *ResponseHostList) Get(b []byte) (*ResponseHostList, *ResponseError, error) {
	res := &ResponseHostList{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseUpdateHostGroup ответ на обновление группы хостов
func NewResponseUpdateHostGroup() *ResponseUpdateHostGroup {
	return &ResponseUpdateHostGroup{}
}

// Get получить список группы хостов
func (r *ResponseUpdateHostGroup) Get(b []byte) (*ResponseUpdateHostGroup, *ResponseError, error) {
	res := &ResponseUpdateHostGroup{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseUpdateHost ответ на обновление хоста
func NewResponseUpdateHost() *ResponseUpdateHost {
	return &ResponseUpdateHost{}
}

// Get получить список хостов
func (r *ResponseUpdateHost) Get(b []byte) (*ResponseUpdateHost, *ResponseError, error) {
	res := &ResponseUpdateHost{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result.HostIds) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseDeleteGroupHost ответ на удаленее группы хостов
func NewResponseDeleteGroupHost() *ResponseDeleteGroupHost {
	return &ResponseDeleteGroupHost{}
}

// Get получить список id удаленных групп хостов
func (r *ResponseDeleteGroupHost) Get(b []byte) (*ResponseDeleteGroupHost, *ResponseError, error) {
	res := &ResponseDeleteGroupHost{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result.GroupIds) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}

// NewResponseDeleteHost ответ на удаленее хостов
func NewResponseDeleteHost() *ResponseDeleteHost {
	return &ResponseDeleteHost{}
}

// Get получить список id удаленных хостов
func (r *ResponseDeleteHost) Get(b []byte) (*ResponseDeleteHost, *ResponseError, error) {
	res := &ResponseDeleteHost{}
	resErr := &ResponseError{}

	res, resErr, err := supportingfunctions.ResponseUnmarchal(b, res, resErr)

	//если нет ошибок но ответ попрежнему пустой
	if len(res.Result.HostIds) == 0 {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}
