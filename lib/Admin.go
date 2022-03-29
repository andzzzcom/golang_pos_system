package LIB

import (
    "html/template"
	"net/http"
	"fmt"
    "golang.org/x/crypto/bcrypt"
	"strings"
)

func ListAdmin(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

    res2 := []Admin{}
	selDB, err := db.Query("SELECT * FROM admins WHERE status!=? ORDER BY id_admin DESC", -1)
	for selDB.Next() {
		var adm Admin
		
		err = selDB.Scan(&adm.Id_admin, &adm.Email, &adm.Password, &adm.Avatar, &adm.Name, &adm.Phone, &adm.Address, &adm.Born_place, &adm.Gender, &adm.Status, &adm.Creator, &adm.Created_date, &adm.Last_updated_date)
		if err != nil{
			fmt.Println(err)
			return
		}
        res2 = append(res2, adm)
	}

    var tmpl = template.Must(template.ParseFiles(
		"views/admin/list.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	var err2 = tmpl.ExecuteTemplate(w, "adminList", res2)
	if err2 != nil {
		fmt.Println(err2)
        return
	}

    defer db.Close()
}

func AddAdmin(w http.ResponseWriter, r *http.Request){
	var data map[string]int

	var tmpl = template.Must(template.ParseFiles(
		"views/admin/add.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "adminAdd", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func AddAdminPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	if r.Method=="POST"{
		name		:= r.FormValue("name")
		email		:= r.FormValue("email")
		born_place	:= r.FormValue("born_place")
		address		:= r.FormValue("address")
		phone		:= r.FormValue("phone")
		gender		:= r.FormValue("gender")
		status		:= r.FormValue("status")

		password		:= "12345678"
		password2, err	:= bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			panic(err.Error())
		}
		
		insForm, err := db.Prepare("INSERT INTO admins(email, password, avatar, name, phone, address, born_place, gender, status, creator) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		if err != nil{
			panic(err.Error())
		}

		insForm.Exec(email, password2, "1.png", name, phone, address, born_place, gender, status, 1)
	}
	defer db.Close()

	http.Redirect(w, r, "/admin", 301)
}

func EditAdmin(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var adm Admin
	selDB, err2 := db.Query("SELECT * FROM admins WHERE id_admin=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&adm.Id_admin, &adm.Email, &adm.Password, &adm.Avatar, &adm.Name, &adm.Phone, &adm.Address, &adm.Born_place, &adm.Gender, &adm.Status, &adm.Creator, &adm.Created_date, &adm.Last_updated_date)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/admin/edit.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3  = tmpl.ExecuteTemplate(w, "adminEdit", adm)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func EditAdminPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	var idAdmin string
	if r.Method=="POST"{
		id			:= r.FormValue("id_admin")
		idAdmin 	= id

		name		:= r.FormValue("name")
		email		:= r.FormValue("email")
		born_place	:= r.FormValue("born_place")
		address		:= r.FormValue("address")
		phone		:= r.FormValue("phone")
		gender		:= r.FormValue("gender")
		status		:= r.FormValue("status")

		updForm, err := db.Prepare("UPDATE admins SET name=?, email=?, born_place=?, address=?, phone=?, gender=?, status=?, status=? WHERE id_admin=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(name, email, born_place, address, phone, gender, status, id)
	}
	defer db.Close()
	
	http.Redirect(w, r, "/admin/edit/"+idAdmin, 301)
}

func DeleteAdmin(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var adm Admin
	selDB, err2 := db.Query("SELECT * FROM admins WHERE id_admin=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&adm.Id_admin, &adm.Email, &adm.Password, &adm.Avatar, &adm.Name, &adm.Phone, &adm.Address, &adm.Born_place, &adm.Gender, &adm.Status, &adm.Creator, &adm.Created_date, &adm.Last_updated_date)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/admin/delete.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3 = tmpl.ExecuteTemplate(w, "adminDelete", adm)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func DeleteAdminPost(w http.ResponseWriter, r *http.Request){
	
    db := DBConnect()

	if r.Method=="POST"{
		id		:= r.FormValue("id_admin")
		updForm, err := db.Prepare("UPDATE admins SET status=? WHERE id_admin=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(-1, id)
	}
	defer db.Close()

	http.Redirect(w, r, "/admin", 301)
}

func GetAdmin() []Admin{
    db := DBConnect()

    res2 := []Admin{}
	selDB, err := db.Query("SELECT * FROM admins WHERE status!=? ORDER BY id_admin DESC", -1)
	for selDB.Next() {
		var adm Admin
		
		err = selDB.Scan(&adm.Id_admin, &adm.Email, &adm.Password, &adm.Avatar, &adm.Name, &adm.Phone, &adm.Address, &adm.Born_place, &adm.Gender, &adm.Status, &adm.Creator, &adm.Created_date, &adm.Last_updated_date)
		if err != nil{
			panic(err.Error())
		}
        res2 = append(res2, adm)
	}
    defer db.Close()

	return res2
}