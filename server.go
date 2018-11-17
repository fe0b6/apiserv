package apiserv

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Начитаем слушать порт
func listen(port int) {

	if initParams.ParseRequest == nil {
		http.HandleFunc("/", parseRequest)
	} else {
		http.HandleFunc("/", initParams.ParseRequest)
	}

	log.Fatalln("[fatal]", http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

// Разбираем запрос
func parseRequest(w http.ResponseWriter, r *http.Request) {

	// Если сервер завершает работу
	if exited {
		w.WriteHeader(503)
		w.Write([]byte(http.StatusText(503)))
		return
	}

	// Отмечаем что начался новый запрос
	wg.Add(1)
	// По завершению запроса отмечаем что он закончился
	defer wg.Done()

	o := &Obj{R: r, W: w, TimeStart: time.Now()}
	initParams.Route(o)
}

// SendAnswer - функция отправки ответа
func (wo *Obj) SendAnswer() {
	// Если ничего не надо делать
	if wo.Ans.Exited {
		return
	}

	// Если нужно вернуть код
	if wo.Ans.Code != 0 {
		wo.sendCode()
		return
	}

	// Добавляем куку если надо
	if wo.Ans.Cookie != "" {
		cookie := http.Cookie{
			Name:     initParams.Cookie.Name,
			Domain:   initParams.Cookie.Domain,
			Path:     initParams.Cookie.Path,
			Value:    wo.Ans.Cookie,
			MaxAge:   initParams.Cookie.Time,
			HttpOnly: true,
			Secure:   initParams.Cookie.Secure,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(wo.W, &cookie)
	}

	// Если переадресация
	if wo.Ans.Redirect != "" {
		wo.W.Header().Add("Expires", "Thu, 01 Jan 1970 00:00:01 GMT")
		http.Redirect(wo.W, wo.R, wo.Ans.Redirect, 301)
		return
	}

	// Добавляем csp
	if len(wo.Ans.CspMap) > 0 {
		csp := getCsp(wo.Ans.CspMap)
		if csp != "" {
			wo.W.Header().Add("Content-Security-Policy", csp)
		}
	} else if initParams.Csp != "" {
		wo.W.Header().Add("Content-Security-Policy", initParams.Csp)
	}

	// Если надо - добавим время выполнения запроса
	if wo.ServerTiming || initParams.ServerTiming {
		t := (time.Now().UnixNano() - wo.TimeStart.UnixNano()) / int64(time.Millisecond)
		wo.W.Header().Add("Server-Timing", "exec;dur="+strconv.FormatInt(t, 10))
	}

	var js []byte
	var err error
	// Формируем json
	if len(wo.Ans.Path) > 0 { // Собираем данные в нужный вид если надо
		var o map[string]interface{}
		o = map[string]interface{}{wo.Ans.Path[len(wo.Ans.Path)-1]: wo.Ans.Data}
		for i := len(wo.Ans.Path) - 2; i >= 0; i-- {
			o = map[string]interface{}{wo.Ans.Path[i]: o}
		}
		js, err = json.Marshal(o)
	} else {
		js, err = json.Marshal(wo.Ans.Data)
	}

	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Пишем ответ
	wo.W.Write(js)
}

// Отправляем ответ
func (wo *Obj) sendCode() {
	wo.W.WriteHeader(wo.Ans.Code)
	// Если не 200 то добавляем статус ответа
	if wo.Ans.Code != 200 {
		wo.W.Write([]byte(http.StatusText(wo.Ans.Code)))
	}
}
