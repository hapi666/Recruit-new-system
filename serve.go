package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/dchest/captcha"
	"database/sql"
	"html/template"
	"log"
	"io"
	"net/http"
)



func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t,err:=template.ParseFiles("3.html")
		if err!=nil {
			fmt.Println(err.Error())
		}
		if r.URL.Path != "/login" {
			http.NotFound(w, r)
			return
		}
		d := struct {
			CaptchaId string
		}{
			captcha.New(),
		}
		if err := t.Execute(w, &d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		fmt.Println("username", r.Form["username"])
		fmt.Println("phonenumber", r.Form["phonenumber"])

		w.Header().Set("Content-Type", "text/template; charset=utf-8")

		db, err := sql.Open("mysql", "root:123456@tcp(119.29.21.123:3306)/test?charset=utf8")
		if err!=nil {
			fmt.Println(err.Error())
		}

		//插入数据
		stmt, err := db.Prepare("INSERT userinfo SET username=?,phonenumber=?,sex=?,class=?,code=?,introduce=?")
		if err!=nil {
			fmt.Println(err.Error())
		}
		r.ParseForm()
		res, err := stmt.Exec(r.FormValue("username"), r.FormValue("phonenumber"),r.FormValue("sex"),r.FormValue("class"),r.FormValue("code"),r.FormValue("introduce"))
		if err!=nil {
			fmt.Println(err.Error())
		}
		if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
			io.WriteString(w, "验证码错误，报名失败！\n")
		} else {
			io.WriteString(w, "报名成功!\n")
		}
		fmt.Println(res)

	}
}

func main() {
	http.HandleFunc("/login", login)         //设置访问的路由
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	err := http.ListenAndServe(":7070", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
