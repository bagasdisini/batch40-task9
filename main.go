package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/detail-project/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/add-project", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", editProject).Methods("GET")
	route.HandleFunc("/update-project/{index}", updateProject).Methods("POST")

	fmt.Println("Server berjalan di port 8080")

	http.ListenAndServe("localhost:8080", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/index.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT name, description, image FROM tb_blog")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.ProjectName, &each.Description, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)
	}
	resData := map[string]interface{}{
		"Project": result,
	}

	template.Execute(w, resData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/contact.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	template.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/detail-project.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	var BlogDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if i == index {
			BlogDetail = Project{
				ProjectName: data.ProjectName,
				Description: data.Description,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				NodeJs:      data.NodeJs,
				ReactJs:     data.ReactJs,
				VueJs:       data.VueJs,
				TypeScript:  data.TypeScript,
				Duration:    data.Duration,
				Image:       data.Image,
			}
		}
	}

	data := map[string]interface{}{
		"BlogDetail": BlogDetail,
	}

	template.Execute(w, data)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var template, error = template.ParseFiles("views/add-project.html")

	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}

	template.Execute(w, nil)
}

type Project struct {
	Id          string
	ProjectName string
	Description string
	StartDate   string
	EndDate     string
	NodeJs      string
	ReactJs     string
	VueJs       string
	TypeScript  string
	Duration    string
	Image       string
}

var dataProject = []Project{}

func addBlog(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	var duration string
	var projectName = r.PostForm.Get("projectName")
	var deskripsi = r.PostForm.Get("deskripsi")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var node = r.PostForm.Get("node")
	var vuejs = r.PostForm.Get("vuejs")
	var react = r.PostForm.Get("react")
	var js = r.PostForm.Get("js")

	var layout = "2006-01-02"
	var startDateParse, _ = time.Parse(layout, startDate)
	var endDateParse, _ = time.Parse(layout, endDate)
	var startDateConvert = startDateParse.Format("02 January 2006")
	var endDateConvert = endDateParse.Format("02 January 2006")

	var hours = endDateParse.Sub(startDateParse).Hours()
	var days = hours / 24
	var weeks = math.Round(days / 7)
	var months = math.Round(days / 30)
	var years = math.Round(days / 365)

	if days >= 1 && days <= 6 {
		duration = strconv.Itoa(int(days)) + " day(s)"
	} else if days >= 7 && days <= 29 {
		duration = strconv.Itoa(int(weeks)) + " week(s)"
	} else if days >= 30 && days <= 364 {
		duration = strconv.Itoa(int(months)) + " month(s)"
	} else if days >= 365 {
		duration = strconv.Itoa(int(years)) + " year(s)"
	}

	var dataBlog = Project{
		ProjectName: projectName,
		Description: deskripsi,
		StartDate:   startDateConvert,
		EndDate:     endDateConvert,
		NodeJs:      node,
		ReactJs:     react,
		VueJs:       vuejs,
		TypeScript:  js,
		Duration:    duration,
	}

	dataProject = append(dataProject, dataBlog)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var editProject = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, project := range dataProject {
		if index == i {
			editProject = Project{
				ProjectName: project.ProjectName,
				Description: project.Description,
				StartDate:   project.StartDate,
				EndDate:     project.EndDate,
				NodeJs:      project.NodeJs,
				ReactJs:     project.ReactJs,
				VueJs:       project.VueJs,
				TypeScript:  project.TypeScript,
				Duration:    project.Duration,
			}
		}

	}

	data := map[string]interface{}{
		"EditProject": editProject,
	}

	tmpl.Execute(w, data)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var duration string
	var projectName = r.PostForm.Get("projectName")
	var deskripsi = r.PostForm.Get("deskripsi")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var node = r.PostForm.Get("node")
	var vuejs = r.PostForm.Get("vuejs")
	var react = r.PostForm.Get("react")
	var js = r.PostForm.Get("js")

	var layout = "2006-01-02"
	var startDateParse, _ = time.Parse(layout, startDate)
	var endDateParse, _ = time.Parse(layout, endDate)
	var startDateConvert = startDateParse.Format("02 January 2006")
	var endDateConvert = endDateParse.Format("02 January 2006")

	var hours = endDateParse.Sub(startDateParse).Hours()
	var days = hours / 24
	var weeks = math.Round(days / 7)
	var months = math.Round(days / 30)
	var years = math.Round(days / 365)

	if days >= 1 && days <= 6 {
		duration = strconv.Itoa(int(days)) + " day(s)"
	} else if days >= 7 && days <= 29 {
		duration = strconv.Itoa(int(weeks)) + " week(s)"
	} else if days >= 30 && days <= 364 {
		duration = strconv.Itoa(int(months)) + " month(s)"
	} else if days >= 365 {
		duration = strconv.Itoa(int(years)) + " year(s)"
	}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject[index] = Project{
		ProjectName: projectName,
		Description: deskripsi,
		StartDate:   startDateConvert,
		EndDate:     endDateConvert,
		NodeJs:      node,
		ReactJs:     react,
		VueJs:       vuejs,
		TypeScript:  js,
		Duration:    duration,
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
