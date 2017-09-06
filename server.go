package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var name string

type user struct {
	Name    string
	Phone   string
	Iclass  string
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if err := t.Execute(w, nil); err != nil {
			fmt.Println(err.Error())
		}
	}
}
func login2(w http.ResponseWriter, r *http.Request) {

	//请求的是登陆数据，那么执行登陆的逻辑判断
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("username", r.Form["username"])
		fmt.Println("phonenumber", r.Form["phonenumber"])

		json, err := ioutil.ReadFile("config_example.json")

		jsonStr := string(json)

		if err != nil {
			fmt.Println(err)
		}

		pd := gjson.Get(jsonStr, "password")

		ppp := pd.String()

		ur := gjson.Get(jsonStr, "user")

		uuu := ur.String()

		dk := gjson.Get(jsonStr, "port")

		ppt := dk.String()

		db, err := sql.Open("mysql", uuu+":"+ppp+"@tcp(119.29.21.123:"+ppt+")/test?charset=utf8")
		if err != nil {
			fmt.Println(err.Error())
		}
		//查找数据
		find, err := db.Query("SELECT id FROM userinfo WHERE name=? AND phone=?", r.FormValue("name"), r.FormValue("phone"))
		user_id := 0
		for find.Next() {
			find.Scan(&user_id)
		}
		if user_id != 0 {
			//更新数据
			up, err := db.Exec("UPDATE userinfo SET iclass=?, message=? WHERE id=?", r.FormValue("iclass"), r.FormValue("message"), user_id)
			fmt.Println(up)
			if err != nil {
				fmt.Println(err)
			}
			r.ParseForm()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			//插入数据
			stmt, err := db.Prepare("INSERT userinfo SET name=?,phone=?,iclass=?,message=?")

			if err != nil {
				fmt.Println(err.Error())
			}
			r.ParseForm()
			res, err := stmt.Exec(r.FormValue("name"), r.FormValue("phone"), r.FormValue("iclass"), r.FormValue("message"))
			if err != nil {
				fmt.Println(err.Error())
				w.Write([]byte("fail"))
			} else {
				w.Write([]byte("success"))
			}
			fmt.Println(res)
		}
		post, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(post)) //将post强制转换成string

		w.WriteHeader(http.StatusOK)

	}
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", login) //设置访问的路由
	http.HandleFunc("/login2", login2)
	files := http.FileServer(http.Dir("Public"))
	rtr.PathPrefix("/").Handler(http.StripPrefix("/", files))
	http.Handle("/", rtr)

	err := http.ListenAndServe("", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
