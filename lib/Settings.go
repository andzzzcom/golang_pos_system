package LIB

import (
    "html/template"
	"net/http"
)

func EditSettingsGeneral(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	var set Settings
	selDB, err2 := db.Query("SELECT * FROM settings WHERE id_setting=?", 2)
	for selDB.Next() {
		err2 = selDB.Scan(&set.Id_setting, &set.Title_web, &set.Subtitle_web, &set.Favicon_web, &set.Logo_web, &set.Email, &set.Status, &set.Creator, &set.Created_date, &set.Last_updated_date)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/settings/general.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3  = tmpl.ExecuteTemplate(w, "settingsEdit", set)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func EditSettingsGeneralPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	if r.Method=="POST"{
		titleWeb	:= r.FormValue("title_web")
		subtitleWeb	:= r.FormValue("subtitle_web")
		email		:= r.FormValue("email")

		updForm, err := db.Prepare("UPDATE settings SET title_web=?, subtitle_web=?, email=? WHERE id_setting=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(titleWeb, subtitleWeb, email, 2)
	}
	defer db.Close()
	
	http.Redirect(w, r, "/settings/general", 301)
}