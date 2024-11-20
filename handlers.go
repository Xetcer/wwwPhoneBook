package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":1234"

/*
defaultHandler Это обработчик по умолчанию, который обслуживает все запросы, не совпадающие ни с одним из других обработчиков.
*/
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

/*
deleteHandler -  функция-обработчик для пути /delete, которая начинается с разделения URL с целью получения нужной информации.
*/
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// получить телефон
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	telephone := paramStr[2]
	err := deleteEntry(telephone)
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}
	Body := telephone + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

/*
list(), которая используется в пути /list, не может завершиться сбоем.
Следовательно, http.StatusOK всегда возвращается при обработке /list.

	Однако иногда возвращаемое значение list() может оказаться пустым.
*/
func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := list()
	fmt.Fprintf(w, "%s", Body)
}

/*
Здесь мы определяем функцию обработчика для URL /status.
Он просто возвращает информацию об общем количестве записей в нашей телефонной книге.
Его можно использовать для проверки того, что веб-сервис работает нормально.
*/
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("total entries: %d\n", len(data))
	fmt.Fprintf(w, "%s", Body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// разделяем URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}
	name := paramStr[2]
	surename := paramStr[3]
	tel := paramStr[4]

	t := strings.ReplaceAll(tel, "-", "")
	if !matchTel(t) {
		fmt.Println("Not a valid telephone number:", tel)
		return
	}

	temp := &Entry{Name: name, Surname: surename, Tel: t}
	err := insert(temp)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		Body := "Failed to add record\n"
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New record added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}
	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

/*
search() проверяет, существует ли данная запись в телефонной книге, и действует соответствующим
Обновление приложения телефонной книги образом.
*/
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// получить значение search из URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}
	var Body string
	telephone := paramStr[2]
	t := search(telephone)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not be found: " + telephone + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = t.Name + " " + t.Surname + " " + t.Tel + "\n"
	}
	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}
