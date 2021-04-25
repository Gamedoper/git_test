package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//!create a struct product and use this model for json

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

//? create slice of struct []Product into product List

var productList []Product

//? initialize product with products Json

func init() {
	productsJson := `[
		{
			"productId": 1,
			"manufacturer": "Johns-Jenkins",
			"sku": "p5z343vdS",
			"pricePerUnit": "497.45",
			"quantityOnHand": 9703,
			"productName": "sticky note"
		  },
		  {
			"productId": 2,
			"manufacturer": "Hessel, Schimmel and Feeney",
			"sku": "i7v300kmx",
			"pricePerUnit": "282.29",
			"quantityOnHand": 9217,
			"productName": "leg warmers"
		  },
		  {
			"productId": 3,
			"manufacturer": "Swaniawski, Bartoletti and Bruen",
			"sku": "q0L657ys7",
			"pricePerUnit": "436.26",
			"quantityOnHand": 5905,
			"productName": "lamp shade"
		  },
		  {
			"productId": 4,
			"manufacturer": "Realiance, Jio",
			"sku": "q0L65a7ys7",
			"pricePerUnit": "436.26",
			"quantityOnHand": 590,
			"productName": "Dongle"
		  }
	]`

	//? we are going to use unmarshal Json data poiting to interface slice product list

	err := json.Unmarshal([]byte(productsJson), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

//* get Id of product and increment it

func getNextId() int {
	highestID := -1
	for _, product := range productList {
		log.Println(product)
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

//* find product by id function

func findProductById(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == productID {
			return &product, i
		}
	}
	return nil, 0
}

//? product handler with switch method get and list out all products

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductID = getNextId()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//* find single product and id

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//! find product by id function is above
	product, listItemIndex := findProductById(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updatedProduct
		productList[listItemIndex] = *product
		w.WriteHeader((http.StatusOK))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	// http.Handle("/", &fooHandler{message: "hello world"}) //! handle we need to pass object
	// http.HandleFunc("/bar", barHandler)                   //! we can just pass the method
	http.HandleFunc("/products/", productHandler)
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":5000", nil) //* listen and serve port

}
