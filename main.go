package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type HomeStruct struct {
	Phones      []Phone
	Message     string
	Search      string
	ResultCount int
}

type Phone struct {
	ID          int
	Name        string
	Price       string
	Image       string
	Brand       string
	Capacity    string
	ScreenSize  string
	Version     string
	FrontCamera string
	BackCamera  string
	Description string
}

type AdminData struct {
	TotalOrders      int
	Orders           []Order
	Phones           []Phone
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
var highest = 9
var phones = []Phone{
	{
		ID:          1,
		Name:        "iPhone 16 Pro",
		Price:       "₦1,850,000",
		Image:       "../static/image/coinview-app-h7a6g0ua6LM-unsplash.jpg",
		Brand:       "Apple",
		Capacity:    "256GB",
		ScreenSize:  "6.7 inches",
		Version:     "iPhone 16 Pro",
		FrontCamera: "12MP",
		Description: "A premium Apple smartphone.",
	},
	{
		ID:          2,
		Name:        "Samsung Galaxy S26",
		Price:       "₦1,650,000",
		Image:       "/static/image/igor-omilaev-X4S-G_Q9U9g-unsplash.jpg",
		Brand:       "Sony",
		Capacity:    "352GB",
		ScreenSize:  "2.5 inches",
		Version:     "Samsung Galaxy S26",
		FrontCamera: "14MP",
		Description: "A premium Samsung smartphone.",
	},
	{
		ID:          3,
		Name:        "Redmi Note 15",
		Price:       "₦450,000",
		Image:       "/static/image/neil-soni-6wdRuK7bVTE-unsplash.jpg",
		Brand:       "Apple",
		Capacity:    "352GB",
		ScreenSize:  "2.5 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "1MP",
		Description: "A premium Redmi smartphone.",
	},
	{
		ID:          4,
		Name:        "Redmi Note 15",
		Price:       "₦450,000",
		Image:       "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand:       "Sony",
		Capacity:    "352GB",
		ScreenSize:  "2.5 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "16MP",
		Description: "A premium Redmi smartphone.",
	},
	{
		ID:          5,
		Name:        "Redmi Note 15",
		Price:       "₦450,00",
		Image:       "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand:       "Google",
		Capacity:    "352GB",
		ScreenSize:  "3.5 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "15MP",
		Description: "A premium Redmi smartphone.",
	},
	{
		ID:          6,
		Name:        "Redmi Note 15",
		Price:       "₦640,000",
		Image:       "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand:       "Xiaomi",
		Capacity:    "352GB",
		ScreenSize:  "2.5 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "14MP",
		Description: "A premium Redmi smartphone.",
	},
	{
		ID:          7,
		Name:        "Redmi Note 15",
		Price:       "₦420,000",
		Image:       "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand:       "Samsung",
		Capacity:    "352GB",
		ScreenSize:  "2.6 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "15MP",
		Description: "A premium Redmi smartphone.",
	},
	{
		ID:          8,
		Name:        "Redmi Note 15",
		Price:       "₦450,000",
		Image:       "/static/image/shiwa-id-Uae7ouMw91A-unsplash.jpg",
		Brand:       "Apple",
		Capacity:    "352GB",
		ScreenSize:  "2.2 inches",
		Version:     "Redmi Note 15",
		FrontCamera: "18MP",
		Description: "A premium Redmi smartphone.",
	},
}

type CartItem struct {
	Phone    Phone
	Quantity int
}

var Cart []CartItem

type CartTotal struct {
	CartItems []CartItem
	Total     int
	ItemCount int
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
		Phones:           phones,
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
	phoneID, _ := strconv.Atoi(id)

	for _, phone := range phones {

		if phone.ID == phoneID {

			for i := 0; i < len(Cart); i++ {
				
				if Cart[i].Phone.ID == phoneID {

					
					Cart[i].Quantity++
					
					http.Redirect(w, r, "/cart", http.StatusSeeOther)
					return
					// tmpl, err := template.ParseFiles("templates/checkout.html")
					// if err != nil {
					// 	http.Error(w, err.Error(), 500)
					// 	return
					// }

					// tmpl.Execute(w, phone)

					// return
				}
			}

			newItem := CartItem{
				Phone:    phone,
				Quantity: 1,
			}
			
			Cart = append(Cart, newItem)
			
			http.Redirect(w, r, "/cart", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "Phone not found", http.StatusNotFound)
}

// func HomeHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles("templates/index.html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	tmpl.Execute(w, phones)
// }

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

func AddingPhoneHandler(w http.ResponseWriter, r *http.Request) {

	error := r.ParseMultipartForm(10 << 20)

	if error != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file inside static/image/
	filePath := "static/image/" + header.Filename

	savedImage, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}
	defer savedImage.Close()

	// Copy uploaded file into static/image/
	_, err = io.Copy(savedImage, file)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}

	AddingPhone := Phone{
		ID:    highest,
		Name:  r.FormValue("name"),
		Price: r.FormValue("price"),

		Brand: r.FormValue("brand"),
		Image: "/" + filePath,
	}

	// fmt.Println(AddingPhone.Image)
	highest++

	for i := 0; i < len(phones); i++ {
		if phones[i].Name == AddingPhone.Name {
			fmt.Fprintln(w, "Phone already exists.")
			return
		}
	}

	phones = append(phones, AddingPhone)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeletePhoneHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("DeletePhoneHandler is running")
	DeletePhone := r.FormValue("delete")

	convertetoint, _ := strconv.Atoi(DeletePhone)

	for i := 0; i < len(phones); i++ {
		if phones[i].ID == convertetoint {
			phones = append(phones[:i], phones[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditingHandler(w http.ResponseWriter, r *http.Request) {

	Editing := r.FormValue("edit")
	convertInt, _ := strconv.Atoi(Editing)

	var Found Phone
	for i := 0; i < len(phones); i++ {
		if phones[i].ID == convertInt {
			Found = phones[i]
			templ, _ := template.ParseFiles("templates/edit.html")
			templ.Execute(w, Found)
			return
		}

	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func SaveEditedPhoneHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ID := r.FormValue("id")

	converInt, _ := strconv.Atoi(ID)

	for i := 0; i < len(phones); i++ {
		if phones[i].ID == converInt {
			phones[i].Name = r.FormValue("name")
			phones[i].Price = r.FormValue("price")
			phones[i].Brand = r.FormValue("brand")
			phones[i].Image = r.FormValue("image")

			break
		}
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)

}

func IndexFunctionHandler(w http.ResponseWriter, r *http.Request) {

	templ, _ := template.ParseFiles("templates/index.html")

	searchText := strings.TrimSpace(r.FormValue("search"))
	sorting := r.FormValue("sort")

	changeToLower := strings.ToLower(searchText)
	var foundphone []Phone

	for i := 0; i < len(phones); i++ {
		if strings.Contains(strings.ToLower(phones[i].Name), changeToLower) ||
			strings.Contains(strings.ToLower(phones[i].Brand), changeToLower) || strings.Contains(phones[i].Price, changeToLower) {
			foundphone = append(foundphone, phones[i])
		}
	}

	switch sorting {
	case "price_low":
		for i := 0; i < len(foundphone); i++ {

			for j := i + 1; j < len(foundphone); j++ {

				cleanPrice1 := foundphone[i].Price
				cleanPrice1 = strings.ReplaceAll(cleanPrice1, ",", "")
				cleanPrice1 = strings.ReplaceAll(cleanPrice1, "₦", "")

				cleanPrice2 := foundphone[j].Price
				cleanPrice2 = strings.ReplaceAll(cleanPrice2, ",", "")
				cleanPrice2 = strings.ReplaceAll(cleanPrice2, "₦", "")

				price1, _ := strconv.Atoi(cleanPrice1)
				price2, _ := strconv.Atoi(cleanPrice2)
				if price1 > price2 {

					temp := foundphone[i]
					foundphone[i] = foundphone[j]
					foundphone[j] = temp
				}

			}

		}
	case "price_high":
		for i := 0; i < len(foundphone); i++ {
			for j := i + 1; j < len(foundphone); j++ {

				cleanPrice1 := foundphone[i].Price
				cleanPrice1 = strings.ReplaceAll(cleanPrice1, ",", "")
				cleanPrice1 = strings.ReplaceAll(cleanPrice1, "₦", "")

				cleanPrice2 := foundphone[j].Price
				cleanPrice2 = strings.ReplaceAll(cleanPrice2, ",", "")
				cleanPrice2 = strings.ReplaceAll(cleanPrice2, "₦", "")

				price1, _ := strconv.Atoi(cleanPrice1)
				price2, _ := strconv.Atoi(cleanPrice2)
				if price1 < price2 {
					temp := foundphone[i]
					foundphone[i] = foundphone[j]
					foundphone[j] = temp
				}
			}
		}
	case "name_az":
		for i := 0; i < len(foundphone); i++ {
			for j := i + 1; j < len(foundphone); j++ {

				Name1 := strings.ToLower(foundphone[i].Name)
				Name2 := strings.ToLower(foundphone[j].Name)
				if Name1 > Name2 {
					temp := foundphone[i]
					foundphone[i] = foundphone[j]
					foundphone[j] = temp
				}

			}
		}
	case "name_za":
		for i := 0; i < len(foundphone); i++ {
			for j := i + 1; j < len(foundphone); j++ {
				Name1 := strings.ToLower(foundphone[i].Name)
				Name2 := strings.ToLower(foundphone[j].Name)
				if Name1 < Name2 {
					temp := foundphone[i]
					foundphone[i] = foundphone[j]
					foundphone[j] = temp
				}
			}
		}
	}

	if len(foundphone) == 0 {
		// templ, _ := template.ParseFiles("templates/index.html")

		HomePageStruct := HomeStruct{
			Phones:  phones,
			Message: "No phones found",
			Search:  searchText,
		}

		templ.Execute(w, HomePageStruct)
		return
	}

	HomePageStruct := HomeStruct{
		Phones:      foundphone,
		Message:     "",
		ResultCount: len(foundphone),
	}

	templ.Execute(w, HomePageStruct)

	// templ, _ := template.ParseFiles("templates/index.html")
	// templ.Execute(w, foundphone)

}

func PhoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	phoneID, _ := strconv.Atoi(id)
	selectedPhone := Phone{}

	for i := 0; i < len(phones); i++ {
		if phones[i].ID == phoneID {

			selectedPhone = phones[i]
			break
		}
	}
	templ, _ := template.ParseFiles("templates/phone.html")
	templ.Execute(w, selectedPhone)
}

func CartHandler(w http.ResponseWriter, r *http.Request) {

	total := 0

	for i := 0; i < len(Cart); i++ {

		cleanPrice1 := Cart[i].Phone.Price
		cleanPrice1 = strings.ReplaceAll(cleanPrice1, ",", "")
		cleanPrice1 = strings.ReplaceAll(cleanPrice1, "₦", "")

		price, _ := strconv.Atoi(cleanPrice1)

		// total = total + price
		total += price * Cart[i].Quantity

	}

	cartData := CartTotal{
		CartItems: Cart,
		Total:     total,
		ItemCount: len(Cart),
	}

	templ, _ := template.ParseFiles("templates/cart.html")
	templ.Execute(w, cartData)
}

func RemoveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := r.FormValue("phone_id")
	ChangeStrInt, _ := strconv.Atoi(id)

	for i := 0; i < len(Cart); i++ {
		if Cart[i].Phone.ID == ChangeStrInt {
			Cart = append(Cart[:i], Cart[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func main() {

	http.HandleFunc("/remove", RemoveHandler)
	http.HandleFunc("/cart", CartHandler)
	http.HandleFunc("/phone", PhoneHandler)
	http.HandleFunc("/", IndexFunctionHandler)
	http.HandleFunc("/savephone", SaveEditedPhoneHandler)

	http.HandleFunc("/edit", EditingHandler)

	http.HandleFunc("/deletephone", DeletePhoneHandler)
	http.HandleFunc("/addphone", AddingPhoneHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/login", AddminLogin)

	http.HandleFunc("/track", TrackFunction)

	http.HandleFunc("/admin", AdminFunction)

	http.HandleFunc("/update", updateFunction)

	http.HandleFunc("/submit", SubmitFuntion)
	http.HandleFunc("/delete", deleteFunction)

	http.HandleFunc("/buy", BuyHandler)

	// http.HandleFunc("/", HomeHandler)
	http.Handle("/images/",
		http.StripPrefix("/images/",
			http.FileServer(http.Dir("./images"))))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/",
		http.StripPrefix("/static/", fs))

	fmt.Println("Server running on http://localhost:8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}
