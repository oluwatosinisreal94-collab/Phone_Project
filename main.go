package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Phone struct {
	ID    int
	Name  string
	Price string
	Image string
	Brand string
}

type AdminData struct {
	TotalOrders      int
	Orders           []Order
	NewOrders        int
	ProcessingOrders int
	CompletedOrders  int
}

type Order struct {
	ID          int
	Customer    string
	PhoneNumber string
	Address     string

	ProductName  string
	ProductPrice string
	Status       string
}

var CustomerInformation = []Order{}

var idCounter = 1
var phones = []Phone{
	{
		ID:    1,
		Name:  "iPhone 16 Pro",
		Price: "₦1,850,000",
		Image: "/static/image/coinview-app-h7a6g0ua6LM-unsplash.jpg",
		Brand: "Apple",
	},
	{
		ID:    2,
		Name:  "Samsung Galaxy S26",
		Price: "₦1,650,000",
		Image: "/static/image/igor-omilaev-X4S-G_Q9U9g-unsplash.jpg",
		Brand: "Sony",
	},
	{
		ID:    3,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/neil-soni-6wdRuK7bVTE-unsplash.jpg",
		Brand: "Apple",
	},
	{
		ID:    4,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand: "Sony",
	},
	{
		ID:    5,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand: "Google",
	},
	{
		ID:    6,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand: "Xiaomi",
	},
	{
		ID:    7,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand: "Samsung",
	},
	{
		ID:    8,
		Name:  "Redmi Note 15",
		Price: "₦450,000",
		Image: "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand: "Apple",
	},
}

func AdminFunction(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("logged_in")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if cookie.Value != "yes" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	TotalOrders := len(CustomerInformation)

	NewOrders := 0
	ProcessingOrders := 0
	CompletedOrders := 0

	for i := 0; i < len(CustomerInformation); i++ {

		switch CustomerInformation[i].Status {
		case "New":
			NewOrders++
		case "Processing":
			ProcessingOrders++
		case "Completed":
			CompletedOrders++
		}
	}

	OrderShow := CustomerInformation

	Search := r.FormValue("search")

	if Search != "" {

		FoundOrder := []Order{}

		for i := 0; i < len(CustomerInformation); i++ {

			if CustomerInformation[i].Customer == Search {

				FoundOrder = append(FoundOrder, CustomerInformation[i])

			}
		}

		OrderShow = FoundOrder
	}

	data := AdminData{
		TotalOrders:      TotalOrders,
		Orders:           OrderShow,
		NewOrders:        NewOrders,
		ProcessingOrders: ProcessingOrders,
		CompletedOrders:  CompletedOrders,
	}

	templ, _ := template.ParseFiles("templates/admin.html")
	templ.Execute(w, data)
}

func IsValidPhone(phone string) bool {

	if len(phone) != 11 {
		return false
	}

	for _, char := range phone {

		if char < '0' || char > '9' {
			return false
		}

	}

	return true
}

func SubmitFuntion(w http.ResponseWriter, r *http.Request) {

	CustomerOrder := Order{
		ID:          idCounter,
		Customer:    r.FormValue("customer"),
		PhoneNumber: r.FormValue("phone_number"),
		Address:     r.FormValue("address"),

		ProductName:  r.FormValue("product_name"),
		ProductPrice: r.FormValue("product_price"),
		Status:       "New",
	}

	if !IsValidPhone(CustomerOrder.PhoneNumber) {
		fmt.Fprint(w, "Invalid phone number. Use only numbers.")
		return
	}

	idCounter++
	CustomerInformation = append(CustomerInformation, CustomerOrder)

	templ, _ := template.ParseFiles("templates/successful.html")
	templ.Execute(w, CustomerOrder)

}

func BuyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", 400)
		return
	}

	id := r.FormValue("phone_id")

	for _, phone := range phones {

		if fmt.Sprint(phone.ID) == id {

			tmpl, err := template.ParseFiles("templates/checkout.html")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			tmpl.Execute(w, phone)

			return
		}
	}

	http.Error(w, "Phone not found", 404)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, phones)
}

func deleteFunction(w http.ResponseWriter, r *http.Request) {
	Delete := r.FormValue("delete")
	changingtoint, _ := strconv.Atoi(Delete)

	for i := 0; i < len(CustomerInformation); i++ {
		if CustomerInformation[i].ID == changingtoint {
			CustomerInformation = append(CustomerInformation[:i], CustomerInformation[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func updateFunction(w http.ResponseWriter, r *http.Request) {

	Update := r.FormValue("update")

	Updatetoint, _ := strconv.Atoi(Update)

	for i := 0; i < len(CustomerInformation); i++ {

		if CustomerInformation[i].ID == Updatetoint {

			switch CustomerInformation[i].Status {
			case "New":

				CustomerInformation[i].Status = "Processing"

			case "Processing":

				CustomerInformation[i].Status = "Completed"

			}
			break
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func TrackFunction(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templ, _ := template.ParseFiles("templates/track.html")
		templ.Execute(w, nil)
		return
	}

	UserReturn := r.FormValue("phone_number")

	// strconvtoInt, _ := strconv.Atoi(UserReturn)

	FoundOrder := Order{}
	found := false

	for i := 0; i < len(CustomerInformation); i++ {

		if CustomerInformation[i].PhoneNumber == UserReturn {

			FoundOrder = CustomerInformation[i]
			found = true
			break

		}

	}
	if !found {
		fmt.Fprintln(w, "Order does not Exist", FoundOrder)
		return
	}
	templ, _ := template.ParseFiles("templates/track.html")
	templ.Execute(w, FoundOrder)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		return
	}

	cookie := http.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func AddminLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templ, _ := template.ParseFiles("templates/login.html")
		templ.Execute(w, nil)
		return
	}

	Admin := r.FormValue("adminname")
	AdminPassword := r.FormValue("adminpassword")

	Name := "Admin"
	Password := 1234

	PasswordConverter, _ := strconv.Atoi(AdminPassword)

	if Admin == Name && PasswordConverter == Password {

		cookie := http.Cookie{
			Name:    "logged_in",
			Value:   "yes",
			Expires: time.Now().Add(1 * time.Hour),
		}

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)

	} else if Admin != Name || PasswordConverter != Password {
		fmt.Fprintln(w, "invalid UserName And Password")
		return
	}

}

func main() {

	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/login", AddminLogin)

	http.HandleFunc("/track", TrackFunction)

	http.HandleFunc("/admin", AdminFunction)

	http.HandleFunc("/update", updateFunction)

	http.HandleFunc("/submit", SubmitFuntion)
	http.HandleFunc("/delete", deleteFunction)

	http.HandleFunc("/buy", BuyHandler)

	http.HandleFunc("/", HomeHandler)
	http.Handle("/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("./images"))))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/",
		http.StripPrefix("/static/", fs))

	fmt.Println("Server running on http://localhost:8084")

	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		fmt.Println(err)
	}
}
