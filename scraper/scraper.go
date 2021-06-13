package main

import (
		"fmt"
		"log"
		"net/http"
		"time"
		"encoding/json"
		"bytes"
		"regexp"
		"strings"
		"strconv"
		"github.com/gorilla/mux"
		"github.com/PuerkitoBio/goquery"	
	)

//------ Type declarations ----------//
type ProductDetails struct {
	Name			string	`json:"name,omitempty"`
	ImageURL		string	`json:"imageURL,omitempty"`
	Description		string	`json:"description,omitempty"`
	Price			string	`json:"price,omitempty"`
	TotalReviews	int		`json:"totalReviews,omitempty"`
}

type Product struct {
	productURL		string			`json:"productURL"`
	Product			ProductDetails	`json:"product,omitempty"`
}

type Status struct {
	ID 				string	`json:"ID,omitempty"`
	Matched			int		`json:"Matched,omitempty"`
    Modified		int		`json:"Modified,omitempty"`
}
//------- end type declarations ------//


//---------- Helper Functions --------//

func getName(document *goquery.Document) string {
	var name string
	h1 := document.Find("h1#title").First()
	name = h1.Find("span#productTitle").Text()
	name = strings.TrimSpace(name)

	if name != "" {
		return name
	} else {
		return "Name Not Found!"
	}
}

func getImageURL(document *goquery.Document) string {
	var imageURL string
	document.Find("div#imgTagWrapperId").First().Each(func(i int, div *goquery.Selection) {
		str, _ := div.Find("img").Attr("data-a-dynamic-image")
		pattern, _:= regexp.Compile("https:\\/\\/.*?.jpg")
		img := pattern.FindAllString(str, -1)
		if len(img) > 1 {
			imageURL = img[len(img)-1]
		}
	})

	if imageURL != "" {
		return imageURL
	} else {
		return "Image URL Not Found!"
	}
}

func getDescription(document *goquery.Document) string {
	var description string
	document.Find("div#feature-bullets").First().Find("li").Each(func(i int, li *goquery.Selection) {
		if i != 0 {
			description += strings.TrimSpace(li.Find("span.a-list-item").Text()) + ". "
		}
	})

	if description != "" {
		return description
	} else {
		return "Description Not Found!"
	}
}

func getPrice(document *goquery.Document) string {
	var price string
	pattern, _:= regexp.Compile("(\\$[,0-9]*(\\.)([0-9])+)")

	price = document.Find("span#priceblock_ourprice").First().Text()

	if !pattern.Match([]byte(price)) {
		str := document.Find("ul.a-unordered-list").First().Find("li#edition_0").First().Find("span.a-size-mini").Text()
		temp := pattern.FindAllString(str, 1)
		if len(temp) > 1 {
			price = temp[0]
		}
	}

	if !pattern.Match([]byte(price)) {
		html, _ := document.Html()
		price = pattern.FindString(string(html))
	}

	if price != "" {
		return price
	} else {
		return "Price Not Found!"
	}
}

func getTotalReviews(document *goquery.Document) int {
	var totalReviews int
	temp_str := document.Find("span#acrCustomerReviewText").First().Text()
	temp_str = strings.ReplaceAll(temp_str, ",", "")
	temp_str = strings.Split(temp_str, " ")[0]
	totalReviews, _ = strconv.Atoi(temp_str)

	return totalReviews
}

func scrape_data(url string) Product {
	client := &http.Client{
        Timeout: 30 * time.Second,
	}
	log.Println("Requesting URL : ",url)
    response, err := client.Get(url)
    if err != nil {
        log.Fatal("GET err (product url): ", err)
    }
	defer response.Body.Close()

    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
		log.Fatal("goquery err: ", err)
	}

	product_details := ProductDetails{
		Name:			getName(document),
		ImageURL:		getImageURL(document),
		Description:	getDescription(document),
		Price:			getPrice(document),
		TotalReviews:	getTotalReviews(document),
	}
	product := Product{
		productURL: 	url,
		Product:		product_details,
	}
	return product
}
//------- end Helper Functions -------//


//-------------- Views ---------------//
func post_scrapeurl_handler(w http.ResponseWriter, request *http.Request){
	request.ParseForm()
	data := Product{}
	data.productURL = request.Form.Get("productURL")


	data = scrape_data(data.productURL)
	log.Println("Scraped successfully!")

	func_data, err := json.Marshal(data)
	if err != nil{
		log.Fatal("Marshal Error : ",err)
	}

	log.Println("JSONified data!")	
	log.Println("Calling dbapi")

	url := "http://dbapi:5001/dbapi"
    requestObject, err := http.NewRequest("POST", url, bytes.NewBuffer(func_data))
    requestObject.Header.Set("content-type", "application/json")

    client := &http.Client{}
    response, err := client.Do(requestObject)
    if err != nil {
        log.Fatal("Response error : ", err)
	}
	log.Println("Response recieved from dbapi")
    defer response.Body.Close()

	var status Status
	_ = json.NewDecoder(response.Body).Decode(&status)

	if status.Matched == 0 {
		fmt.Fprintf(w, "Product details scraped and stored in database")
	}


}

func get_scrapeurl_handler(response http.ResponseWriter, request *http.Request){
	fmt.Fprintf(response, "Make a POST request!")
}
//------- end Views -------//




func main(){
	fmt.Println("Go application service started!\nPlease make a POST request with x-www-form-urlencoded data\nkey : productURL , value : <product url>")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scrapeapi", post_scrapeurl_handler).Methods("POST")
	router.HandleFunc("/scrapeapi", get_scrapeurl_handler).Methods("GET")
	http.ListenAndServe(":5000", router)
}