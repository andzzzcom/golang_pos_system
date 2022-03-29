package main

import(
    "html/template"
	"net/http"
	"github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
	"github.com/kataras/go-sessions/v3"
	"path"
	LIB "2/lib"
)

type Product struct{
	Id_product			int	
	Name				string
	Price				string
	Category			int	
	Thumbnail			string
	Created_date		string
	Last_updated_date	string
	Creator				int
	Status				int	
}

type Category struct{
	Id_category			int	
	Name				string
	Created_date		string
	Last_updated_date	string
	Creator				int
	Status				int	
}

func main(){
	r := mux.NewRouter()
	
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", loginHandlerPost).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	//invoice
	r.HandleFunc("/invoice", LIB.ListInvoice).Methods("GET")
	r.HandleFunc("/invoice/add", LIB.AddInvoice).Methods("GET")
	r.HandleFunc("/invoice/add_act", LIB.AddInvoicePost).Methods("POST")
	r.HandleFunc("/invoice/edit/{id}", LIB.EditInvoice).Methods("GET")
	r.HandleFunc("/invoice/edit_act", LIB.EditInvoicePost).Methods("POST")
	r.HandleFunc("/invoice/delete/{id}", LIB.DeleteInvoice).Methods("GET")
	r.HandleFunc("/invoice/delete_act", LIB.DeleteInvoicePost).Methods("POST")

	//product
	r.HandleFunc("/product", LIB.ListProduct).Methods("GET")
	r.HandleFunc("/product/add", LIB.AddProduct).Methods("GET")
	r.HandleFunc("/product/add_act", LIB.AddProductPost).Methods("POST")
	r.HandleFunc("/product/edit/{id}", LIB.EditProduct).Methods("GET")
	r.HandleFunc("/product/edit_act", LIB.EditProductPost).Methods("POST")
	r.HandleFunc("/product/delete/{id}", LIB.DeleteProduct).Methods("GET")
	r.HandleFunc("/product/delete_act", LIB.DeleteProductPost).Methods("POST")

	//category
	r.HandleFunc("/category", LIB.ListCategory).Methods("GET")
	r.HandleFunc("/category/add", LIB.AddCategory).Methods("GET")
	r.HandleFunc("/category/add_act", LIB.AddCategoryPost).Methods("POST")
	r.HandleFunc("/category/edit/{id}", LIB.EditCategory).Methods("GET")
	r.HandleFunc("/category/edit_act", LIB.EditCategoryPost).Methods("POST")
	r.HandleFunc("/category/delete/{id}", LIB.DeleteCategory).Methods("GET")
	r.HandleFunc("/category/delete_act", LIB.DeleteCategoryPost).Methods("POST")

	//admin
	r.HandleFunc("/admin", LIB.ListAdmin).Methods("GET")
	r.HandleFunc("/admin/add", LIB.AddAdmin).Methods("GET")
	r.HandleFunc("/admin/add_act", LIB.AddAdminPost).Methods("POST")
	r.HandleFunc("/admin/edit/{id}", LIB.EditAdmin).Methods("GET")
	r.HandleFunc("/admin/edit_act", LIB.EditAdminPost).Methods("POST")
	r.HandleFunc("/admin/delete/{id}", LIB.DeleteAdmin).Methods("GET")
	r.HandleFunc("/admin/delete_act", LIB.DeleteAdminPost).Methods("POST")
	
	//settings
	r.HandleFunc("/settings/general", LIB.EditSettingsGeneral).Methods("GET")
	r.HandleFunc("/settings/general", LIB.EditSettingsGeneralPost).Methods("POST")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
  
	http.ListenAndServe(":8000", r)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	var data = map[string]int{
		"s":1,
		"t":2,
	}
    fp := path.Join("views", "login.html")
    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request){
    db := LIB.DBConnect()

	if r.Method=="POST"{
		email		:= r.FormValue("email")
		password	:= r.FormValue("password")
		
		var adm LIB.Admin
		selDB, err2 := db.Query("SELECT * FROM admins WHERE email=?", email)
		for selDB.Next() {
			err2 = selDB.Scan(&adm.Id_admin, &adm.Email, &adm.Password, &adm.Avatar, &adm.Name, &adm.Phone, &adm.Address, &adm.Born_place, &adm.Gender, &adm.Status, &adm.Creator, &adm.Created_date, &adm.Last_updated_date)
			if err2 != nil{
				panic(err2.Error())
			}
		}
		//deskripsi dan compare password
		var check_pass = bcrypt.CompareHashAndPassword([]byte(adm.Password), []byte(password))

		if check_pass == nil {
			//login success
			
			session := sessions.Start(w, r)
			session.Set("name", adm.Name)
			session.Set("email", adm.Email)
			http.Redirect(w, r, "/", 302)
		} else {
			//login failed
			http.Redirect(w, r, "/login", 302)
		}
	}
	defer db.Close()
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/login", 302)
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	invoice	:= LIB.GetInvoice()
	product	:= LIB.GetProduct()
	category:= LIB.GetCategory()
	admin	:= LIB.GetAdmin()

	session := sessions.Start(w, r)
	if len(session.GetString("email")) == 0 {
		http.Redirect(w, r, "/login", 301)
	}

	var data = map[string]interface{}{
		"name":session.GetString("name"),
		"email":session.GetString("email"),
		"invoice":invoice,
		"product":product,
		"category":category,
		"admin":admin,
	}
    var tmpl = template.Must(template.ParseFiles(
		"views/index.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	var err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
