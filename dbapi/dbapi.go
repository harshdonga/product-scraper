package main

import (
	"time"
	"log"
	"context"
	"net/http"
	"encoding/json"
	"hash/fnv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductDetails struct {
	Name			string	`json:"name,omitempty"`
	ImageURL		string	`json:"imageURL,omitempty"`
	Description		string	`json:"description,omitempty"`
	Price			string	`json:"price,omitempty"`
	TotalReviews	int		`json:"totalReviews,omitempty"`
}

type Product struct {
	ID				uint32				`json:"_id,omitempty" bson:"_id,omitempty"`
	productURL		string				`json:"productURL,omitempty"`
	Product			ProductDetails		`json:"product,omitempty"`
	LastUpdate		time.Time			`json:"last_update, omitempty" bson:"last_update, omitempty"`

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}


var client *mongo.Client

func FindOrAdd(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var new_doc, existing_doc Product
	_ = json.NewDecoder(request.Body).Decode(&new_doc)
	collection := client.Database("productdb").Collection("sellerappcollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	collection.FindOne(ctx, bson.M{"_id":hash(new_doc.Product.Name)}).Decode(&existing_doc)
	new_doc.LastUpdate = time.Now()
	new_doc.ID = hash(new_doc.Product.Name)
	// log.Println("existing doc : ",existing_doc)
	// log.Println("new doc : ", new_doc)
	if existing_doc.ID == 0 {
		log.Println("Creating new document!")
		result, _ := collection.InsertOne(ctx, new_doc)
		json.NewEncoder(response).Encode(result)
	} else {
		log.Println("Updating existing document")
		result, _ := collection.UpdateOne(	ctx,
								bson.M{"_id": hash(new_doc.Product.Name)},
								bson.D{
									primitive.E{
										Key: "$set",
										Value: bson.D{
												primitive.E{
													Key: "product",
													Value: new_doc.Product,
												},
										},
									},
								},
							)
		json.NewEncoder(response).Encode(result)
	}
	
}


func GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var docs []Product
	collection := client.Database("productdb").Collection("sellerappcollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc Product
		cursor.Decode(&doc)
		docs = append(docs, doc)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	
	json.NewEncoder(response).Encode(docs)
}


func main()  {
	localhost := "mongodb://database:27017"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(localhost)
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/dbapi", FindOrAdd).Methods("POST")
	router.HandleFunc("/dbapi", GetAll).Methods("GET")
	http.ListenAndServe(":5001", router)
}
