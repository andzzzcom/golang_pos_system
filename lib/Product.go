package LIB

import (
    "html/template"
	"net/http"
	"fmt"
	"github.com/dustin/go-humanize"
	"strconv"
	"strings"
)

type Invoice struct{
	Id_invoice			int	
	Invoice_code		string
	Subtotal_price		string
	Tax_price			string	
	Total_price			string
	Status				int	
	Creator				int
	Created_date		string
	Last_updated_date	string
}

type InvoiceDetail struct{
	Id_invoice_detail	int	
	Id_invoice			int	
	Id_product			int
	Status				int	
	Created_date		string
	Last_updated_date	string
}

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

type Admin struct{
	Id_admin			int	
	Email				string
	Password			string
	Avatar				string	
	Name				string
	Phone				string
	Address				string
	Born_place			string
	Gender				int
	Status				int	
	Creator				int
	Created_date		string
	Last_updated_date	string
}

type Settings struct{
	Id_setting			int	
	Title_web			string
	Subtitle_web		string
	Favicon_web			string	
	Logo_web			string
	Email				string
	Status				int	
	Creator				int
	Created_date		string
	Last_updated_date	string
}

func ListProduct(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

    res2 := []Product{}
	selDB, err := db.Query("SELECT * FROM products WHERE status!=? ORDER BY id_product DESC", -1)
	for selDB.Next() {
		var prod Product

		err = selDB.Scan(&prod.Id_product, &prod.Name, &prod.Price, &prod.Category, &prod.Thumbnail, &prod.Created_date, &prod.Last_updated_date, &prod.Creator, &prod.Status)
		if err != nil{
			fmt.Println(err)
			return
		}

		prod.Price = convertCurrency(prod.Price)
        res2 = append(res2, prod)
	}

    var tmpl = template.Must(template.ParseFiles(
		"views/product/list.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	var err2 = tmpl.ExecuteTemplate(w, "productList", res2)
	if err2 != nil {
		fmt.Println(err2)
        return
	}

    defer db.Close()
}

func AddProduct(w http.ResponseWriter, r *http.Request){
	categories	:= GetCategory()
	
	var tmpl = template.Must(template.ParseFiles(
		"views/product/add.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "productAdd", categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func AddProductPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	if r.Method=="POST"{
		name	:= r.FormValue("name")
		price	:= r.FormValue("price")
		status	:= r.FormValue("status")

		insForm, err := db.Prepare("INSERT INTO products(name, price, category, thumbnail, creator, status) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil{
			panic(err.Error())
		}

		insForm.Exec(name, price, 1, "1.png", 1, status)
		fmt.Println("Sukses!")
	}
	defer db.Close()

	http.Redirect(w, r, "/product", 301)
}

func EditProduct(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	categories	:= GetCategory()

	id := strings.Split(r.URL.Path, "/")[3]

	var prod Product
	selDB, err2 := db.Query("SELECT * FROM products WHERE id_product=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&prod.Id_product, &prod.Name, &prod.Price, &prod.Category, &prod.Thumbnail, &prod.Created_date, &prod.Last_updated_date, &prod.Creator, &prod.Status)
		if err2 != nil{
			fmt.Println(err2)
			return
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/product/edit.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	data := map[string]interface{}{
        "prod": prod,
        "cat": categories,
    }

	var err3  = tmpl.ExecuteTemplate(w, "productEdit", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func EditProductPost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	var idProduct string
	if r.Method=="POST"{
		id		:= r.FormValue("id_product")
		idProduct = id

		name	:= r.FormValue("name")
		price	:= r.FormValue("price")
		status	:= r.FormValue("status")

		updForm, err := db.Prepare("UPDATE products SET name=?, price=?, category=?, thumbnail=?, creator=?, status=? WHERE id_product=?")
		if err != nil{
			panic(err.Error())
		}

		updForm.Exec(name, price, 1, "2.png", 1, status, id)
		fmt.Println("Sukses!")
	}
	defer db.Close()

	http.Redirect(w, r, "/product/edit/"+idProduct, 301)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var prod Product
	selDB, err2 := db.Query("SELECT * FROM products WHERE id_product=?", id)
	for selDB.Next() {
		err2 = selDB.Scan(&prod.Id_product, &prod.Name, &prod.Price, &prod.Category, &prod.Thumbnail, &prod.Created_date, &prod.Last_updated_date, &prod.Creator, &prod.Status)
		if err2 != nil{
			fmt.Println(err2)
			return
		}
		fmt.Println(prod)
	}
	fmt.Println(prod)

	categories	:= GetCategory()
	data := map[string]interface{}{
        "prod": prod,
        "cat": categories,
    }

	var tmpl = template.Must(template.ParseFiles(
		"views/product/delete.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err3 = tmpl.ExecuteTemplate(w, "productDelete", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func DeleteProductPost(w http.ResponseWriter, r *http.Request){
	
    db := DBConnect()

	if r.Method=="POST"{
		id		:= r.FormValue("id_product")
		updForm, err := db.Prepare("UPDATE products SET status=? WHERE id_product=?")
		if err != nil{
			panic(err.Error())
		}

		updForm.Exec(-1, id)
		fmt.Println("Sukses!")
	}
	defer db.Close()

	http.Redirect(w, r, "/product", 301)
}

func GetProduct() []Product{
    db := DBConnect()

    res2 := []Product{}
	selDB, err := db.Query("SELECT * FROM products WHERE status!=? ORDER BY id_product DESC", -1)
	for selDB.Next() {
		var prod Product

		err = selDB.Scan(&prod.Id_product, &prod.Name, &prod.Price, &prod.Category, &prod.Thumbnail, &prod.Created_date, &prod.Last_updated_date, &prod.Creator, &prod.Status)
		if err != nil{
			panic(err.Error())
		}
        res2 = append(res2, prod)
	}
    defer db.Close()
	return res2
}

func convertCurrency(price string) string{
	str, _ := strconv.ParseFloat(price, 32)
	result := strings.ReplaceAll(humanize.Commaf(str), ",", ".")
	return result
}