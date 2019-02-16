package apiserv

import (
	"context"
	"log"
	"strings"
	"sync"
)

var (
	wg         sync.WaitGroup
	exited     bool
	initParams Param
)

// Init это функция инициализации
func Init(p Param) (exitChan chan bool) {
	// Устанавливаем параметры по умолчанию
	setDefaultParams(p)

	// Собираем Content-Security-Policy
	setCsp()

	// Канал для оповещения о выходе
	exitChan = make(chan bool)

	// Начинаем слушать http-порт
	go listen(initParams.Port)

	go waitExit(exitChan)

	return
}

// Ждем сигнал о выходе
func waitExit(exitChan chan bool) {
	_ = <-exitChan

	//exited = true

	log.Println("[info]", "Завершаем работу api сервера")

	// Ждем пока все запросы завершатся
	//wg.Wait()

	if err := httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v\n", err)
	}

	log.Println("[info]", "Работа api сервера завершена корректно")

	exitChan <- true
}

// Устанавливаем параметры по умолчанию
func setDefaultParams(p Param) {
	initParams = p

	if initParams.Cookie.Path == "" {
		initParams.Cookie.Path = "/"
	}
}

// Собираем Content-Security-Policy
func setCsp() {
	if initParams.CspMap == nil {
		return
	}

	csp := []string{}
	for k, v := range initParams.CspMap {
		csp = append(csp, k+" "+v)
	}

	initParams.Csp = strings.Join(csp, "; ")
}

func getCsp(h map[string]string) string {
	if initParams.CspMap == nil {
		return ""
	}

	csp := []string{}
	for k, v := range initParams.CspMap {
		if _, ok := h[k]; ok {
			v = v + " " + h[k]
		}
		csp = append(csp, k+" "+v)
	}

	return strings.Join(csp, "; ")
}
