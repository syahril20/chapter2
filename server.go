package main

import (
	"context"
	"example/hello/connection"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
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
	e.Logger.Fatal(e.Start("localhost:2009"))
}

var node string
var next string
var react string
var typescript string

func index(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id_project, img, title, start, stop, duration, postDate, description, technologies FROM tb_project")
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
	Proj := map[string]interface{}{
		"Project": result,
	}

	fmt.Println(result[1])
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
	layout := ("2006-01-02")
	img := "/public/assets/content.png"
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

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project ( img, title, start, stop, duration, postdate,description, technologies) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", img, title, start, end, duration, postdate, description, technologies)

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
	var result []Project
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

	var data, _ = connection.Conn.Query(context.Background(), "SELECT id_project, img, title, start, stop, postdate, description, technologies FROM tb_project where id_project = $1", id)
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
