package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"sync"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

var lock sync.Mutex

// album represents data about a record album.
type user struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Posts    []string `json:"posts"`
}

type post struct {
	ID               string `json:"id"`
	AuthorId         string `json:"authorid"`
	Caption          string `json:"caption"`
	Image_Url        string `json:"image_url"`
	Posted_Timestamp int64  `json:"posted_timestamp"`
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func postUsers(collection *mongo.Collection) gin.HandlerFunc {
	lock.Lock()
    defer lock.Unlock()
	fn := func(c *gin.Context) {
		var newUser user

		if err := c.BindJSON(&newUser); err != nil {
			return
		}

		hashedPassword := GetMD5Hash(newUser.Password)
		fmt.Println(hashedPassword)

		newUser.Password = hashedPassword
		newUser.Posts = []string{}

		insertResult, err := collection.InsertOne(context.TODO(), newUser)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		c.IndentedJSON(http.StatusCreated, newUser)
	}

	return gin.HandlerFunc(fn)
}

func getUserByID(collection *mongo.Collection) gin.HandlerFunc {
	lock.Lock()
    defer lock.Unlock()
	fn := func(c *gin.Context) {
		id := c.Param("id")
		filter := bson.D{{Key: "id", Value: id}}
		var result user

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		c.IndentedJSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

func postPosts(postCollection *mongo.Collection, usersCollection *mongo.Collection) gin.HandlerFunc {
	lock.Lock()
    defer lock.Unlock()
	fn := func(c *gin.Context) {
		var newPost post
		authorID := newPost.AuthorId

		if err := c.BindJSON(&newPost); err != nil {
			return
		}
		now := time.Now()
		newPost.Posted_Timestamp = now.Unix()
		fmt.Print(authorID)
		userFilter := bson.D{{Key: "id", Value: newPost.AuthorId}}
		
		var result user

		err := usersCollection.FindOne(context.TODO(), userFilter).Decode(&result)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Author not found"})
			return
		}
		insertResult, err := postCollection.InsertOne(context.TODO(), newPost)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		}
		
		filter := bson.D{{Key: "id", Value: authorID}}

		postID := newPost.ID
		update := bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "posts", Value: postID},
			}},
		}
		updateResult, err := usersCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		c.IndentedJSON(http.StatusCreated, newPost)
	}
	return gin.HandlerFunc(fn)
}

func getPostByID(collection *mongo.Collection) gin.HandlerFunc {
	lock.Lock()
    defer lock.Unlock()
	fn := func(c *gin.Context) {
		id := c.Param("id")
		filter := bson.D{{Key: "id", Value: id}}
		var result post

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, result)
	}
	return gin.HandlerFunc(fn)
}

func getPostsOfAuthor(collection *mongo.Collection) gin.HandlerFunc {
	lock.Lock()
    defer lock.Unlock()
	fn := func(c *gin.Context) {

		query := c.Request.URL.Query()
		fmt.Println(reflect.TypeOf(query["page"][0]))
		offset, err := strconv.Atoi(query["page"][0])
		if err != nil {
			log.Fatal(err)
		}
		offset -=1
		offset *= 3
		fmt.Println(offset)
		options := options.Find()
		options.SetLimit(3)
		options.SetSkip(int64(offset))

		authorId := c.Param("authorId")
		filter := bson.D{{Key: "authorid", Value: authorId}}
		var results []*post

		cur, err := collection.Find(context.TODO(), filter, options)
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var elem post
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		cur.Close(context.TODO())

		fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
		if len(results) == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, results)
	}
	return gin.HandlerFunc(fn)
}

func main() {

	clientOptions := options.Client().ApplyURI("mongodb+srv://harshul:harshul@cluster0.4eyuv.mongodb.net/GoDB?retryWrites=true&w=majority")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	usersCollection := client.Database("InstaAPI").Collection("users")
	postsCollection := client.Database("InstaAPI").Collection("posts")

	router := gin.Default()
	router.POST("/users", postUsers(usersCollection))
	router.POST("/posts", postPosts(postsCollection, usersCollection))
	router.GET("/users/:id", getUserByID(usersCollection))
	router.GET("/posts/:id", getPostByID(postsCollection))
	router.GET("/posts/users/:authorId", getPostsOfAuthor(postsCollection))

	router.Run("localhost:8080")
}
