package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/tarantool/go-tarantool"
	"runtime"
	"sync"
	"os/exec"
)


var TARANTOOL_SERVER string = "127.0.0.1:4502"


func main(){

	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	wg.Add(2)

	// Сервер обработки HTTP-запросов
	go func() {

		defer wg.Done()

		http.HandleFunc("/", index_handler)
		http.HandleFunc("/get", get_handler)
		http.HandleFunc("/set", set_handler)
		http.ListenAndServe("localhost:8000", nil)
		fmt.Println("Server started")
	}()

	// Сервер Tarantool
	go func() {

		defer wg.Done()

		cmd := exec.Command("tarantool", "init.lua")
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Unable to start Tarantool", err)
		}

	}()

	wg.Wait()

}

// Заглавная страница
func index_handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,	 "Tarantool DB interface" +
		"\n\nPrint '/get?id=x' to request value with id=x" +
		"\n\nPrint '/set?val=x' to add value with val=x")
}


// Страничка вывода полученных данных
func get_handler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	for key, values := range r.Form {
		if key == "id"{
			id, err :=  strconv.Atoi(values[0])
			if err == nil {
				if id >= 0{
					fmt.Println("GET: id = ", id)
					value := get_value(id)
					fmt.Println(value[0])
					rcv := value[0]
					fmt.Fprintln(w, "Got:", rcv)
				} else {
					fmt.Println("Invalid request")
					fmt.Fprintf(w, "Invalid request")
				}
			} else {
				fmt.Println("Invalid request")
				fmt.Fprintf(w, "Invalid request")
			}
		}
	}
}

// Получение данных из тарантула
func get_value(id int) (val []interface{}){

	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect(TARANTOOL_SERVER, opts)
	res := []interface{}{"Data unavailable"}

	if err != nil {
		fmt.Println("Connection refused:", err)
	} else {
		fmt.Println("Connected to Tarantool")
		resp, err := conn.Call("get_value", []interface{}{id})

		if err != nil {
			fmt.Println("Error", err)
		} else {
			res = resp.Data
			fmt.Println(res[0])
		}
		conn.Close()
		if err == nil {
			fmt.Println("Connection closed")
		}else{
			fmt.Println("Failed to close connection: ", err)
		}
	}

	return res
}

// Страничка результат ввода данных
func set_handler(w http.ResponseWriter, r *http.Request){

	r.ParseForm()
	for key, values := range r.Form {
		if key == "val"{
			value := values[0]
			fmt.Println("SET: val = ", value)
			id := set_value([]interface{} {value})
			fmt.Println(id[0])
			fmt.Fprintln(w, "Added with ID", id[0])
		}
	}
}

// Внесение новых данных в тарантул
func set_value(val []interface{}) (id []interface{}){

	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect(TARANTOOL_SERVER, opts)
	res := []interface{}{"Wrong request"}

	if err != nil {
		fmt.Println("Connection refused:", err)
	} else {
		fmt.Println("Connected to Tarantool")
		resp, err := conn.Call("set_value", val)

		if err != nil {
			fmt.Println("Error", err)
		} else {
			res = resp.Data
		}
		err = conn.Close()
		if err == nil {
			fmt.Println("Connection closed")
		}else{
			fmt.Println("Failed to close connection: ", err)
		}
	}

	return res
}


