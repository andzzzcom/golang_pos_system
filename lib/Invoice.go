package LIB

import (
    "html/template"
	"net/http"
	"fmt"
	"strconv"
	"strings"
)

func ListInvoice(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

    res2 := []Invoice{}
	selDB, err := db.Query("SELECT * FROM invoices WHERE status!=? ORDER BY id_invoice DESC", -1)
	for selDB.Next() {
		var inv Invoice

		err = selDB.Scan(&inv.Id_invoice, &inv.Invoice_code, &inv.Subtotal_price, &inv.Tax_price, &inv.Total_price, &inv.Status, &inv.Creator, &inv.Created_date, &inv.Last_updated_date)
		if err != nil{
			fmt.Println(err)
			return
		}
		
		inv.Subtotal_price = convertCurrency(inv.Subtotal_price)
		inv.Tax_price = convertCurrency(inv.Tax_price)
		inv.Total_price = convertCurrency(inv.Total_price)
        res2 = append(res2, inv)
	}

    var tmpl = template.Must(template.ParseFiles(
		"views/invoice/list.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	var err2 = tmpl.ExecuteTemplate(w, "invoiceList", res2)
	if err2 != nil {
		fmt.Println(err2)
        return
	}

    defer db.Close()
}

func AddInvoice(w http.ResponseWriter, r *http.Request){
	products	:= GetProduct()

	var tmpl = template.Must(template.ParseFiles(
		"views/invoice/add.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "invoiceAdd", products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func AddInvoicePost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	if r.Method=="POST"{
		id_product, err			:= strconv.Atoi(r.FormValue("id_product"))
		code					:= r.FormValue("invoice_code")
		tax_price, errs			:= strconv.Atoi(r.FormValue("tax_price"))
		status					:= r.FormValue("status")
		subtotal_price, err4	:= strconv.Atoi(getPriceProduct(id_product))
		total_price				:= tax_price+subtotal_price
		
		if err != nil{
			panic(err.Error())
		}
		if errs != nil{
			panic(errs.Error())
		}
		if err4 != nil{
			panic(err4.Error())
		}

		insForm, err := db.Prepare("INSERT INTO invoices(invoice_code, subtotal_price, tax_price, total_price, creator, status) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil{
			panic(err.Error())
		}
		res, err5 := insForm.Exec(code, subtotal_price, tax_price, total_price, 1, status)
		if err5 != nil{
			panic(err5.Error())
		}

		lid, errs2 := res.LastInsertId()
		if errs2 != nil{
			panic(errs2.Error())
		}
		
		insForm2, err2 := db.Prepare("INSERT INTO invoices_detail(id_invoice, id_product, status) VALUES(?, ?, ?)")
		if err2 != nil{
			panic(err2.Error())
		}
		insForm2.Exec(lid, id_product, 1)
	}
	defer db.Close()

	http.Redirect(w, r, "/invoice", 301)
}

func EditInvoice(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var inv Invoice
	var invD InvoiceDetail
	selDB, err2 := db.Query("SELECT i.*, ids.id_product FROM invoices i, invoices_detail ids WHERE i.id_invoice=? AND ids.id_invoice=i.id_invoice", id)
	for selDB.Next() {
		err2 = selDB.Scan(&inv.Id_invoice, &inv.Invoice_code, &inv.Subtotal_price, &inv.Tax_price, &inv.Total_price, &inv.Status, &inv.Creator, &inv.Created_date, &inv.Last_updated_date, &invD.Id_product)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/invoice/edit.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))
	
	products	:= GetProduct()
	data := map[string]interface{}{
        "prod": products,
        "inv": inv,
        "invD": invD,
    }

	var err3  = tmpl.ExecuteTemplate(w, "invoiceEdit", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func EditInvoicePost(w http.ResponseWriter, r *http.Request){	
    db := DBConnect()

	var idInvoice string
	if r.Method=="POST"{
		id		:= r.FormValue("id_invoice")
		idInvoice = id

		id_product, err			:= strconv.Atoi(r.FormValue("id_product"))
		code					:= r.FormValue("invoice_code")
		tax_price, errs			:= strconv.Atoi(r.FormValue("tax_price"))
		status					:= r.FormValue("status")
		subtotal_price, err4	:= strconv.Atoi(getPriceProduct(id_product))
		total_price				:= tax_price+subtotal_price

		if err != nil{
			panic(err.Error())
		}
		if errs != nil{
			panic(errs.Error())
		}
		if err4 != nil{
			panic(err4.Error())
		}

		updForm, err := db.Prepare("UPDATE invoices SET invoice_code=?, tax_price=?, subtotal_price=?, total_price=?, status=? WHERE id_invoice=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(code, tax_price, subtotal_price, total_price, status, id)

		updForm2, err2 := db.Prepare("UPDATE invoices_detail SET id_product=? WHERE id_invoice=?")
		if err2 != nil{
			panic(err2.Error())
		}
		updForm2.Exec(id_product, id)
	}
	defer db.Close()

	http.Redirect(w, r, "/invoice/edit/"+idInvoice, 301)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request){
    db := DBConnect()

	id := strings.Split(r.URL.Path, "/")[3]

	var inv Invoice
	var invD InvoiceDetail
	selDB, err2 := db.Query("SELECT i.*, ids.id_product FROM invoices i, invoices_detail ids WHERE i.id_invoice=? AND ids.id_invoice=i.id_invoice", id)
	for selDB.Next() {
		err2 = selDB.Scan(&inv.Id_invoice, &inv.Invoice_code, &inv.Subtotal_price, &inv.Tax_price, &inv.Total_price, &inv.Status, &inv.Creator, &inv.Created_date, &inv.Last_updated_date, &invD.Id_product)
		if err2 != nil{
			panic(err2.Error())
		}
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/invoice/delete.html",
		"views/partials/_header_settings.html",
		"views/partials/_header.html",
		"views/partials/_sidebar.html",
		"views/partials/_other.html",
		"views/partials/_footer.html",
	))

	products	:= GetProduct()
	data := map[string]interface{}{
        "prod": products,
        "inv": inv,
        "invD": invD,
    }

	var err3 = tmpl.ExecuteTemplate(w, "invoiceDelete", data)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	defer db.Close()
}

func DeleteInvoicePost(w http.ResponseWriter, r *http.Request){
	
    db := DBConnect()

	if r.Method=="POST"{
		id		:= r.FormValue("id_invoice")
		updForm, err := db.Prepare("UPDATE invoices SET status=? WHERE id_invoice=?")
		if err != nil{
			panic(err.Error())
		}
		updForm.Exec(-1, id)
	}
	defer db.Close()

	http.Redirect(w, r, "/invoice", 301)
}

func getPriceProduct(idProduct int) string{
    db := DBConnect()

	var prod Product
	selDB, err := db.Query("SELECT * FROM products WHERE status!=? AND id_product=?", -1, idProduct)
	for selDB.Next() {
		err = selDB.Scan(&prod.Id_product, &prod.Name, &prod.Price, &prod.Category, &prod.Thumbnail, &prod.Created_date, &prod.Last_updated_date, &prod.Creator, &prod.Status)
		if err != nil{
			panic(err.Error())
		}
	}
    defer db.Close()
	return prod.Price
}

func GetInvoice() []Invoice{
    db := DBConnect()

    res2 := []Invoice{}
	selDB, err := db.Query("SELECT * FROM invoices WHERE status!=? ORDER BY id_invoice DESC", -1)
	for selDB.Next() {
		var inv Invoice

		err = selDB.Scan(&inv.Id_invoice, &inv.Invoice_code, &inv.Subtotal_price, &inv.Tax_price, &inv.Total_price, &inv.Status, &inv.Creator, &inv.Created_date, &inv.Last_updated_date)
		if err != nil{
			panic(err.Error())
		}
		
		inv.Subtotal_price = convertCurrency(inv.Subtotal_price)
		inv.Tax_price = convertCurrency(inv.Tax_price)
		inv.Total_price = convertCurrency(inv.Total_price)
        res2 = append(res2, inv)
	}
    defer db.Close()

	return res2
}
