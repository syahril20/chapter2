package main

import (
	"context"
	"example/hello/connection"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Id          int
	Img         string
	Title       string
	Start       string
	End         string
	Duration    string
	Postdate    string
	Description string
	React       bool
	Next        bool
	Node        bool
	Typescript  bool
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
	e.GET("/", index)
	e.GET("/myProject", myProject)
	e.GET("/editProject/:id", editProject)
	e.GET("/contactMe", contactMe)
	e.POST("/add-proj", addProj)
	e.POST("/edit-proj/:id", updateProj)
	e.GET("/proj-detail/:id", projDetail)
	e.GET("/delete-proj/:id", deleteProj)

	// Server
	e.Logger.Fatal(e.Start("localhost:2000"))
}

func index(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	var data, _ = connection.Conn.Query(context.Background(), "SELECT id_project, img, title, duration, postDate, description, react, next, node, typescript FROM tb_project")
	var result []Project
	for data.Next() {
		var p = Project{}
		err := data.Scan(&p.Id, &p.Img, &p.Title, &p.Duration, &p.Postdate, &p.Description, &p.React, &p.Next, &p.Node, &p.Typescript)
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

func myProject(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/myProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
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
	img := "/public/assets/content.png"
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	duration := createDuration(start, end)
	postdate := time.Now().Format("02/01/2006")
	description := c.FormValue("description")
	react := c.FormValue("reactBox") == "on"
	next := c.FormValue("nextBox") == "on"
	node := c.FormValue("nodeBox") == "on"
	typescript := c.FormValue("typeBox") == "on"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project ( img, title, start, stop, duration, postdate,description, react, next, node, typescript) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", img, title, start, end, duration, postdate, description, react, next, node, typescript)

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

	var data, _ = connection.Conn.Query(context.Background(), "SELECT id_project, img, title, start, stop,duration,postDate, description, react, next, node, typescript FROM tb_project where id_project = $1", id)
	var result []Project
	// var Prok = Project{}
	for data.Next() {
		var p = Project{}
		err := data.Scan(&p.Id, &p.Img, &p.Title, &p.Start, &p.End, &p.Duration, &p.Postdate, &p.Description, &p.React, &p.Next, &p.Node, &p.Typescript)
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

func updateProj(c echo.Context) error {
	id := c.Param("id")
	img := "/public/assets/content.png"
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	duration := createDuration(start, end)
	postdate := time.Now().Format("02/01/2006")
	description := c.FormValue("description")
	node := c.FormValue("nodeBox") == "on"
	next := c.FormValue("nextBox") == "on"
	react := c.FormValue("reactBox") == "on"
	typescript := c.FormValue("typeBox") == "on"

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project	SET img=$1, title=$2, start=$3, stop=$4, description=$5, react=$6, next=$7, node=$8, typescript=$9, duration=$10, postdate=$11	WHERE id_project=$12;", img, title, start, end, description, react, next, node, typescript, duration, postdate, id)
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

	var data, _ = connection.Conn.Query(context.Background(), "SELECT id_project, img, title, start, stop,duration,postDate, description, react, next, node, typescript FROM tb_project where id_project = $1", id)
	var result []Project
	// var Prok = Project{}
	for data.Next() {
		var p = Project{}
		err := data.Scan(&p.Id, &p.Img, &p.Title, &p.Start, &p.End, &p.Duration, &p.Postdate, &p.Description, &p.React, &p.Next, &p.Node, &p.Typescript)
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
	// dataProj = append(dataProj[:id], dataProj[id+1:]...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
