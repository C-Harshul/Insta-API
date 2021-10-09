// package main

// import(
// 	"testing"
// 	"net/http"
// 	"net/http/httptest"
// 	"context"
// 	"log"
// 	"fmt"

// 	//"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	"github.com/gin-gonic/gin"

// )

// func TestGetUser(t *testing.T) {

// 	clientOptions := options.Client().ApplyURI("mongodb+srv://harshul:harshul@cluster0.4eyuv.mongodb.net/GoDB?retryWrites=true&w=majority")

// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		t.Errorf("Failed to connect to Database")
// 	}

// 	err = client.Ping(context.TODO(), nil)

// 	if err != nil {
// 		t.Errorf("Failed while pinging Database")
// 	}
// 	req, err := http.NewRequest("GET", "/users/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(getUserByID(usersCollection))
// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	expected := `[{"id":1,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb@g.com","phone_number":"0987654321"},{"id":2,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"},{"id":6,"first_name":"FirstNameSample","last_name":"LastNameSample","email_address":"lr@gmail.com","phone_number":"1111111111"}]`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }


// func Router() {

// 	clientOptions := options.Client().ApplyURI("mongodb+srv://harshul:harshul@cluster0.4eyuv.mongodb.net/GoDB?retryWrites=true&w=majority")

// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = client.Ping(context.TODO(), nil)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB!")

// 	usersCollection := client.Database("InstaAPI(Test)").Collection("users")
// 	postsCollection := client.Database("InstaAPI(Test)").Collection("posts")

// 	router := gin.Default()
// 	router.POST("/users", postUsers(usersCollection))
// 	router.POST("/posts", postPosts(postsCollection, usersCollection))
// 	router.GET("/users/:id", getUserByID(usersCollection))
// 	router.GET("/posts/:id", getPostByID(postsCollection))
// 	router.GET("/posts/users/:authorId", getPostsOfAuthor(postsCollection))

// 	router.Run("localhost:8080")
// }