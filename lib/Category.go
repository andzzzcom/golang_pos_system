package LIB

import (
    "html/template"
	"net/http"
	"fmt"
	"strings"
)

func ListCategory(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

    res2 := []Category{}
	selDB, err := db.Query("SELECT * FROM categories WHERE status!=? ORDER BY id_category DESC", -1)
	for selDB.Next() {
		var cat Category

		err = selDB.Scan(&cat.Id_category, &cat.Name, &cat.Created_date, &cat.Last_updated_date, &cat.Creator, &cat.Status)
		if err != nil{
			fmt.Println(err)
			return
		}
        res2 = append(res2, cat)
	}

    var tmpl = template.Must(template.ParseFiles(
		"views/category/list.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	var err2 = tmpl.ExecuteTemplate(w, "categoryList", res2)
	if err2 != nil {
		fmt.Println(err2)
        return
	}

    defer db.Close()
}

func AddCategory(w http.ResponseWriter, r *http.Request){
	
	var data map[string]interface{}

	var tmpl = template.Must(template.ParseFiles(
		"views/category/add.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "categoryAdd", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func AddCategoryPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	if r.Method=="POST"{
		name	:= r.FormValue("name")
		status	:= r.FormValue("status")

		insForm, err := db.Prepare("INSERT INTO categories(name, creator, status) VALUES(?, ?, ?)")
		if err != nil{
			panic(err.Error())
		}
		insForm.Exec(name, 1, status)
	}
	defer db.Close()

	http.Redirect(w, r, "/category", 301)
}

func EditCategory(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var cat Category
	selDB, err2 := db.Query("SELECT * FROM categories WHERE id_category=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&cat.Id_category, &cat.Name, &cat.Created_date, &cat.Last_updated_date, &cat.Creator, &cat.Status)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/category/edit.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3  = tmpl.ExecuteTemplate(w, "categoryEdit", cat)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func EditCategoryPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	var idCategory string
	if r.Method=="POST"{
		id		:= r.FormValue("id_category")
		idCategory = id

		name	:= r.FormValue("name")
		status	:= r.FormValue("status")

		updForm, err := db.Prepare("UPDATE categories SET name=?, status=? WHERE id_category=?")
		if err != nil{
			panic(err.Error())
		}

		updForm.Exec(name, status, id)
		fmt.Println("Sukses!")
	}
	defer db.Close()

	http.Redirect(w, r, "/category/edit/"+idCategory, 301)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var cat Category
	selDB, err2 := db.Query("SELECT * FROM categories WHERE id_category=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&cat.Id_category, &cat.Name, &cat.Created_date, &cat.Last_updated_date, &cat.Creator, &cat.Status)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/category/delete.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3 = tmpl.ExecuteTemplate(w, "categoryDelete", cat)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func DeleteCategoryPost(w http.ResponseWriter, r *http.Request){
	
    db := DBConnect()

	if r.Method=="POST"{
		id		:= r.FormValue("id_category")
		updForm, err := db.Prepare("UPDATE categories SET status=? WHERE id_Category=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(-1, id)
	}
	defer db.Close()

	http.Redirect(w, r, "/category", 301)
}


func GetCategory() []Category{
    db := DBConnect()

    res2 := []Category{}
	selDB, err := db.Query("SELECT * FROM categories WHERE status!=? ORDER BY id_category DESC", -1)
	for selDB.Next() {
		var cat Category

		err = selDB.Scan(&cat.Id_category, &cat.Name, &cat.Created_date, &cat.Last_updated_date, &cat.Creator, &cat.Status)
		if err != nil{
			panic(err.Error())
		}
        res2 = append(res2, cat)
	}
    defer db.Close()

	return res2
}