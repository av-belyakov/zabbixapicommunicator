package connectionjsonrpc

import (
	"encoding/json"
)

// NewResponseCreateHostGroup ответ на создание группы хостов
func NewResponseCreateHostGroup() *ResponseCreateHostGroup {
	return &ResponseCreateHostGroup{}
}

// ResponseCreateHostGroup создание группы хостов
func (r *ResponseCreateHostGroup) Get(b []byte) (*ResponseCreateHostGroup, *ResponseError, error) {
	res := &ResponseCreateHostGroup{}
	resErr := &ResponseError{}

	err := json.Unmarshal(b, res)
	if err != nil {
		if err = json.Unmarshal(b, resErr); err != nil {
			return res, resErr, err
		}

		return res, resErr, nil
	}

	return res, resErr, nil
}

// NewResponseGetHostGroupList ответ на запрос с целью получения списка групп хостов
func NewResponseGetHostGroupList() *ResponseHostGroupList {
	return &ResponseHostGroupList{}
}

// ResponseGetHostGroupList получение списка групп хостов
func (r *ResponseHostGroupList) Get(b []byte) (*ResponseHostGroupList, *ResponseError, error) {
	res := &ResponseHostGroupList{}
	resErr := &ResponseError{}

	err := json.Unmarshal(b, res)
	if err != nil {
		if err = json.Unmarshal(b, resErr); err != nil {
			return res, resErr, err
		}

		return res, resErr, nil
	}

	return res, resErr, nil
}

// NewResponseGetHostList ответ на запрос с целью получения списка хостов
func NewResponseGetHostList() *ResponseHostList {
	return &ResponseHostList{}
}

// ResponseHostList получение списка хостов
func (r *ResponseHostList) Get(b []byte) (*ResponseHostList, *ResponseError, error) {
	res := &ResponseHostList{}
	resErr := &ResponseError{}

	err := json.Unmarshal(b, res)
	if err != nil {
		if err = json.Unmarshal(b, resErr); err != nil {
			return res, resErr, err
		}

		return res, resErr, nil
	}

	return res, resErr, nil
}
