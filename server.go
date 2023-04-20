package main

import (
	"html/template"
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

var dataProj = []Project{
	{
		Img:         "/public/assets/content.png",
		Title:       "Phyton Jadi Bahasa Pemrograman Terpopuler di Dunia, Ini Alasannya.",
		Start:       "2023-04-17",
		End:         "2025-05-17",
		Duration:    "3 Bulan",
		Postdate:    "16/04/2023",
		Description: "PARBOABOA - Tiobe merilis peringkat 10 bahasa pemograman paling populer di dunia untuk Oktober 2021. Dalam laporan bertajuk 'Tiobe Programming Community index', Phyton dinobatkan sebagai bahasa programming terpopuler di dunia saat ini.",
		React:       true,
		Next:        true,
		Node:        false,
		Typescript:  true,
	},
}

func main() {
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
	e.Logger.Fatal(e.Start("localhost:1234"))
}

func index(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	Proj := map[string]interface{}{
		"Project": dataProj,
	}

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

func addProj(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	img := c.FormValue("img")
	title := c.FormValue("title")
	start := c.FormValue("start")
	end := c.FormValue("end")
	duration := c.FormValue("duration")
	description := c.FormValue("description")

	var addProj = Project{
		Id:          id,
		Img:         img,
		Title:       title,
		Start:       start,
		End:         end,
		Duration:    duration,
		Postdate:    time.Now().String(),
		Description: description,
		React:       false,
		Next:        false,
		Node:        false,
		Typescript:  false,
	}

	if c.Request().Method == "POST" {
		if c.Request().FormValue("reactBox") == "on" {
			addProj.React = true
		}
		if c.Request().FormValue("nextBox") == "on" {
			addProj.Next = true
		}
		if c.Request().FormValue("nodeBox") == "on" {
			addProj.Node = true
		}
		if c.Request().FormValue("typeBox") == "on" {
			addProj.Typescript = true
		}
	}

	dataProj = append(dataProj, addProj)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var projData = Project{}
	for index, data := range dataProj {
		if id == index {
			projData = Project{
				Id:          id,
				Img:         data.Img,
				Title:       data.Title,
				Start:       data.Start,
				End:         data.End,
				Duration:    data.Duration,
				Postdate:    data.Postdate,
				Description: data.Description,
				React:       data.React,
				Next:        data.Next,
				Node:        data.Node,
				Typescript:  data.Typescript,
			}
		}
	}

	data := map[string]interface{}{
		"Project": projData,
	}

	return tmpl.Execute(c.Response(), data)
}

func updateProj(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	img := c.FormValue("img")
	title := c.FormValue("title")
	// start := c.FormValue("start")
	// end := c.FormValue("end")
	duration := c.FormValue("duration")
	description := c.FormValue("description")
	// node := c.FormValue("nodeBox") != ""
	// next := c.FormValue("nextBox") != ""
	// react := c.FormValue("reactBox") != ""
	// typescript := c.FormValue("typeBox") != ""
	// for i, project := range dataProj {
	// 	if project.Id == id {
	// 		dataProj[i].Title = title
	// 		dataProj[i].Start = start
	// 		dataProj[i].End = end
	// 		dataProj[i].Duration = duration
	// 		dataProj[i].Postdate = postDate
	// 		dataProj[i].Description = description
	// 		dataProj[i].Node = node
	// 		dataProj[i].Next = next
	// 		dataProj[i].React = react
	// 		dataProj[i].Typescript = typescript
	// 	}
	// }

	for i := range dataProj {
		if dataProj[i].Id == id {
			// memperbarui data dengan data yang baru
			dataProj[i].Img = img
			dataProj[i].Title = title
			dataProj[i].Duration = duration
			dataProj[i].Description = description

			// mengembalikan tampilan daftar data dengan data yang telah diperbarui
			return c.Render(http.StatusOK, "/", dataProj)
		}
	}
	// var p = Project{
	// 	Id:          id,
	// 	Img:         img,
	// 	Title:       title,
	// 	Start:       start,
	// 	End:         end,
	// 	Duration:    duration,
	// 	Postdate:    time.Now().String(),
	// 	Description: description,
	// 	React:       react,
	// 	Next:        next,
	// 	Node:        node,
	// 	Typescript:  typescript,
	// }

	// return c.JSON(http.StatusOK, p)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/projDetail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var projData = Project{}
	for index, data := range dataProj {
		if id == index {
			projData = Project{
				Img:         data.Img,
				Title:       data.Title,
				Duration:    data.Duration,
				Postdate:    data.Postdate,
				Description: data.Description,
				React:       data.React,
				Next:        data.Next,
				Node:        data.Node,
				Typescript:  data.Typescript,
			}
		}
	}

	data := map[string]interface{}{
		"Project": projData,
	}

	return tmpl.Execute(c.Response(), data)
}

func deleteProj(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProj = append(dataProj[:id], dataProj[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}
