package main

import (
	"context"
	"example/hello/connection"
	"example/hello/middleware"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id           int
	Img          string
	Title        string
	Start        string
	End          string
	Duration     string
	Postdate     string
	Description  string
	Technologies []bool
	UserName     string
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

// var dataProj = []Project{
// 	{
// 		Img:         "/public/assets/content.png",
// 		Title:       "Phyton Jadi Bahasa Pemrograman Terpopuler di Dunia, Ini Alasannya.",
// 		Start:       "2023-04-17",
// 		End:         "2025-05-17",
// 		Duration:    "3 Bulan",
// 		Postdate:    "16/04/2023",
// 		Description: "PARBOABOA - Tiobe merilis peringkat 10 bahasa pemograman paling populer di dunia untuk Oktober 2021. Dalam laporan bertajuk 'Tiobe Programming Community index', Phyton dinobatkan sebagai bahasa programming terpopuler di dunia saat ini.",
// 		React:       true,
// 		Next:        true,
// 		Node:        false,
// 		Typescript:  true,
// 	},
// }

func main() {
	connection.DatabaseConnect()
	e := echo.New()
	e.Static("/public", "public")
	e.Static("/upload", "upload")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))
	e.GET("/", index)
	e.GET("/loginForm", loginForm)
	e.GET("/myProject", myProject)
	e.GET("/editProject/:id", editProject)
	e.GET("/contactMe", contactMe)
	e.GET("/succes", succes)
	e.GET("/proj-detail/:id", projDetail)
	e.GET("/delete-proj/:id", deleteProj)
	e.POST("/add-proj", middleware.UploadFile(addProj))
	e.POST("/edit-proj/:id", middleware.UploadFile(updateProj))
	e.POST("/login", login)
	e.POST("/register", register)
	e.GET("/logout", logout)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusNotFound
		message := "404 Page Not Found"
		if e, ok := err.(*echo.HTTPError); ok {
			code = e.Code
			message = e.Message.(string)
		}
		c.HTML(code, "<h1>"+message+"</h1>")
	}
	// Server
	e.Logger.Fatal(e.Start("localhost:2054"))
}

var node string
var next string
var react string
var typescript string

func index(c echo.Context) error {
	sess, _ := session.Get("session", c)

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data, err := connection.Conn.Query(context.Background(), "select tb_project.id_project, img, title, start,stop, duration, postdate,description, technologies, tb_user.name from tb_project inner join tb_user on tb_user.id_user =  tb_project.id_user order by tb_project.id_project asc")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	result := []Project{}

	for data.Next() {
		var p = Project{}
		var techArray pgtype.VarcharArray

		err := data.Scan(
			&p.Id,
			&p.Img,
			&p.Title,
			&p.Start,
			&p.End,
			&p.Duration,
			&p.Postdate,
			&p.Description,
			&techArray,
			&p.UserName,
		)
		if err != nil {
			log.Fatal(err)
		}

		p.Technologies = make([]bool, len(techArray.Elements))
		for i, el := range techArray.Elements {
			if el.String == "node" {
				p.Technologies[i] = true
			} else if el.String == "next" {
				p.Technologies[i] = true
			} else if el.String == "react" {
				p.Technologies[i] = true
			} else if el.String == "typescript" {
				p.Technologies[i] = true
			} else {
				p.Technologies[i] = false
			}
		}

		result = append(result, p)

	}

	dataS, err := connection.Conn.Query(context.Background(), "select tb_project.id_project, img, title, start,stop, duration, postdate,description, technologies, tb_user.name from tb_project inner join tb_user on tb_user.id_user =  tb_project.id_user where tb_user.id_user = $1 order by tb_project.id_project asc", sess.Values["id"])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	resultS := []Project{}

	for dataS.Next() {
		var p = Project{}
		var techArray pgtype.VarcharArray

		err := dataS.Scan(
			&p.Id,
			&p.Img,
			&p.Title,
			&p.Start,
			&p.End,
			&p.Duration,
			&p.Postdate,
			&p.Description,
			&techArray,
			&p.UserName,
		)
		if err != nil {
			log.Fatal(err)
		}

		p.Technologies = make([]bool, len(techArray.Elements))
		for i, el := range techArray.Elements {
			if el.String == "node" {
				p.Technologies[i] = true
			} else if el.String == "next" {
				p.Technologies[i] = true
			} else if el.String == "react" {
				p.Technologies[i] = true
			} else if el.String == "typescript" {
				p.Technologies[i] = true
			} else {
				p.Technologies[i] = false
			}
		}

		resultS = append(resultS, p)

	}

	Proj := map[string]interface{}{
		"Project":     result,
		"ProjectS":    resultS,
		"alertStatus": sess.Values["alertStatus"],
		"isLogin":     sess.Values["isLogin"],
		"name":        sess.Values["name"],
		"id":          sess.Values["id"],
	}

	delete(sess.Values, "alertStatus")
	sess.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), Proj)
}

func succes(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/success.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func myProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/myProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	sess, _ := session.Get("session", c)
	Proj := map[string]interface{}{
		"isLogin": sess.Values["isLogin"],
	}
	return tmpl.Execute(c.Response(), Proj)
}

func contactMe(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contactMe.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func createDuration(start string, end string) string {
	startTime, _ := time.Parse("2006-01-02", start)
	endTime, _ := time.Parse("2006-01-02", end)

	duration := endTime.Sub(startTime)
	days := int(duration.Hours() / 24)
	months := int(math.Floor(float64(days) / 30))
	years := int(math.Floor(float64(months) / 12))

	if days > 0 && days <= 29 {
		return fmt.Sprintf("%d Hari", days)
	} else if days >= 30 && months <= 12 {
		return fmt.Sprintf("%d Bulan", months)
	} else if months >= 12 {
		return fmt.Sprintf("%d Tahun", years)
	} else if days >= 0 && days <= 24 {
		return "1 Hari"
	}
	return ""
}

func addProj(c echo.Context) error {
	sess, _ := session.Get("session", c)
	layout := ("2006-01-02")
	img := c.Get("dataFile").(string)
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	duration := createDuration(start, end)
	postdate := time.Now().Format(layout)
	description := c.FormValue("description")

	if c.FormValue("nodeBox") == "on" {
		node = "node"
	} else {
		node = "fNode"
	}
	if c.FormValue("nextBox") == "on" {
		next = "next"
	} else {
		next = "fNext"
	}
	if c.FormValue("reactBox") == "on" {
		react = "react"
	} else {
		react = "fReact"
	}
	if c.FormValue("typeBox") == "on" {
		typescript = "typescript"
	} else {
		typescript = "fTypescript"
	}

	technologies := []string{
		node,
		next,
		react,
		typescript,
	}

	fmt.Println(technologies)

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project ( img, title, start, stop, duration, postdate,description, technologies, id_user) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", img, title, start, end, duration, postdate, description, technologies, sess.Values["id"])

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	// Return success response
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var data, _ = connection.Conn.Query(context.Background(), "SELECT id_project, img, title, start, stop,duration,postDate, description, technologies FROM tb_project where id_project = $1", id)
	var ResultD []Project
	// var Prok = Project{}
	for data.Next() {
		var p = Project{}
		var techArray pgtype.VarcharArray
		err := data.Scan(
			&p.Id,
			&p.Img,
			&p.Title,
			&p.Start,
			&p.End,
			&p.Duration,
			&p.Postdate,
			&p.Description,
			&techArray,
		)
		p.Technologies = make([]bool, len(techArray.Elements))
		for i, el := range techArray.Elements {
			if el.String == "node" {
				p.Technologies[i] = true
			} else if el.String == "next" {
				p.Technologies[i] = true
			} else if el.String == "react" {
				p.Technologies[i] = true
			} else if el.String == "typescript" {
				p.Technologies[i] = true
			} else {
				p.Technologies[i] = false
			}
		}
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"massage": err.Error()})
		}
		ResultD = append(ResultD, p)
	}
	sess, _ := session.Get("session", c)
	Proj := map[string]interface{}{
		"Project": ResultD,
		"isLogin": sess.Values["isLogin"],
	}
	fmt.Println(Proj)

	return tmpl.Execute(c.Response(), Proj)
}

func updateProj(c echo.Context) error {
	id := c.Param("id")
	img := c.Get("dataFile").(string)
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	duration := createDuration(start, end)
	postdate := time.Now().Format("2006-01-02")
	description := c.FormValue("description")
	if c.FormValue("nodeBox") == "on" {
		node = "node"
	} else {
		node = "fNode"
	}
	if c.FormValue("nextBox") == "on" {
		next = "next"
	} else {
		next = "fNext"
	}
	if c.FormValue("reactBox") == "on" {
		react = "react"
	} else {
		react = "fReact"
	}
	if c.FormValue("typeBox") == "on" {
		typescript = "typescript"
	} else {
		typescript = "fTypescript"
	}

	technologies := []string{
		node,
		next,
		react,
		typescript,
	}

	fmt.Println(img)
	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project	SET img=$1, title=$2, start=$3, stop=$4, description=$5, duration=$6, postdate=$7, technologies=$8	WHERE id_project=$9;", img, title, start, end, description, duration, postdate, technologies, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	// Return success response
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projDetail(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/projDetail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var data, _ = connection.Conn.Query(context.Background(), "select tb_project.id_project, img, title, start,stop, duration, postdate,description, technologies, tb_user.name from tb_project inner join tb_user on tb_user.id_user =  tb_project.id_user where tb_project.id_project = $1", id)
	var result []Project
	for data.Next() {
		var p = Project{}
		var techArray pgtype.VarcharArray

		err := data.Scan(
			&p.Id,
			&p.Img,
			&p.Title,
			&p.Start,
			&p.End,
			&p.Duration,
			&p.Postdate,
			&p.Description,
			&techArray,
			&p.UserName,
		)

		p.Technologies = make([]bool, len(techArray.Elements))
		for i, el := range techArray.Elements {
			if el.String == "node" {
				p.Technologies[i] = true
			} else if el.String == "next" {
				p.Technologies[i] = true
			} else if el.String == "react" {
				p.Technologies[i] = true
			} else if el.String == "typescript" {
				p.Technologies[i] = true
			} else {
				p.Technologies[i] = false
			}
		}
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"massage": err.Error()})
		}
		result = append(result, p)
	}

	Proj := map[string]interface{}{
		"Project": result,
	}
	fmt.Println(Proj)

	return tmpl.Execute(c.Response(), Proj)
}

func deleteProj(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id_project=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	// Return success response
	return redirectWithMessage(c, "Berhasil tambah data", true, "/loginForm")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status

	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}

func loginForm(c echo.Context) error {
	sess, _ := session.Get("session", c)
	Proj := map[string]interface{}{
		"isLogin": sess.Values["isLogin"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "alertStatus")
	sess.Save(c.Request(), c.Response())
	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), Proj)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("nameLog")
	password := c.FormValue("passLog")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE name=$1", name).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	fmt.Println(user)
	if err != nil {
		return redirectWithMessage(c, "Name Salah !", false, "/loginForm")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	fmt.Println(err)
	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/loginForm")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 jam
	sess.Values["message"] = "Login Success"
	sess.Values["alertStatus"] = true // show alert
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.Id
	sess.Values["isLogin"] = true // access login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	// sess.Values["isLogin"] = nil
	// sess.Values["name"] = nil
	delete(sess.Values, "isLogin")
	delete(sess.Values, "name")
	sess.Save(c.Request(), c.Response())

	// Hapus semua cache yang terkait dengan halaman saat ini
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	return c.Redirect(http.StatusMovedPermanently, "/")
}
