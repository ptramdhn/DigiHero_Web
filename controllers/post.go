package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	Id                   int64  `json:"id"`
	ApplicantsName       string `json:"applicant_name"`
	EventName            string `json:"event_name"`
	Date                 string `json:"date"`
	EventVenues          string `json:"event_venues"`
	RequirementMaterials string `json:"requirement_materials"`
}

var BASE_URL = "http://localhost:9888/api/v1"

func Index(w http.ResponseWriter, r *http.Request) {

	var posts []Post

	response, err := http.Get(BASE_URL + "/applicants")
	if err != nil {
		log.Print(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&posts); err != nil {
		log.Print(err)
	}

	data := map[string]interface{}{
		"posts": posts,
	}

	temp, _ := template.ParseFiles("views/index.html")
	temp.Execute(w, data)

}

func Create(w http.ResponseWriter, r *http.Request) {

	var post Post
	var data map[string]interface{}

	id := r.URL.Query().Get("id")

	fmt.Println(id)

	if id != "" {
		res, err := http.Get(BASE_URL + "/applicants/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&post); err != nil {
			log.Print(err)
		}

		data = map[string]interface{}{
			"post": post,
		}
	}

	temp, _ := template.ParseFiles("views/create.html")
	temp.Execute(w, data)

}

func Store(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id := r.Form.Get("post_id")

	idInt, _ := strconv.ParseInt(id, 10, 64)
	newPost := Post{
		Id:                   idInt,
		ApplicantsName:       r.Form.Get("post_name"),
		EventName:            r.Form.Get("post_event"),
		Date:                 r.Form.Get("post_date"),
		EventVenues:          r.Form.Get("post_venue"),
		RequirementMaterials: r.Form.Get("post_material"),
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request
	var err error

	if id != "" {
		//update
		fmt.Println("Proses update")
		req, err = http.NewRequest(http.MethodPatch, BASE_URL+"/applicants/"+id, buff)
	} else {
		// create
		fmt.Println("Proses create")
		req, err = http.NewRequest(http.MethodPost, BASE_URL+"/applicants", buff)
	}

	if err != nil {
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var postResponse Post

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&postResponse); err != nil {
		log.Print(err)
	}

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)
	fmt.Println(postResponse)

	if res.StatusCode == 201 || res.StatusCode == 200 {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	req, err := http.NewRequest(http.MethodDelete, BASE_URL+"/applicants/"+id, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)

	if res.StatusCode == 200 {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}

}
