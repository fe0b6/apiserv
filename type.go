package apiserv

import (
	"net/http"
	"time"
)

// Param это переменные для инициализации класса
type Param struct {
	Port         int
	Route        func(*Obj)
	Cookie       Cookie
	CspMap       map[string]string
	Csp          string
	ParseRequest func(http.ResponseWriter, *http.Request)
	ServerTiming bool
}

// Cookie - Объект с описание кукисов
type Cookie struct {
	Name   string
	Domain string
	Path   string
	Time   int
	Secure bool
}

// Obj основной объект запроса
type Obj struct {
	W            http.ResponseWriter
	R            *http.Request
	Ans          Answer
	AppendFunc   func(*Obj, map[string]interface{}) map[string]interface{}
	TimeStart    time.Time
	Debug        bool
	ServerTiming bool
}

// Answer объект содержащий ответ
type Answer struct {
	Path     []string
	Redirect string
	Cookie   string
	Data     interface{}
	Exited   bool
	Code     int
	CspMap   map[string]string
}
