package main

import (
	"10x-go/models"
	"context"
	"log"
	"net/http"
	"os"

	// "firebase/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type Firestore struct {
	client *firestore.Client
	ctx    context.Context
}

func main() {

	// Initial firestore database
	f := Firestore{}
	f.Init()

	// Initial gingonic
	router := gin.New()

	// Index page
	router.GET("/", f.heartbeat)

	// Set username password for access resource
	router.Use(gin.BasicAuth(gin.Accounts{
		"admin": os.Getenv("BASIC_AUTH_PASSWORD"),
	}))

	// Route path in application
	applicationRoute := router.Group("/applications")
	{
		applicationRoute.GET("/", f.getAllApplication)
		applicationRoute.POST("/", f.createApplication)
		applicationRoute.PUT("/:id", f.updateApplicaiton)
		applicationRoute.DELETE("/:id", f.deleteApplication)
	}

	// Start server (production use port 5000)
	router.Run(":5000")
}

// Connect firestore database using credential
func (route *Firestore) Init() {
	route.ctx = context.Background()
	sa := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(route.ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	route.client, err = app.Firestore(route.ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// Homepage
func (route *Firestore) heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, 10x GO"})
}

// Return all applications in database as json format
func (route *Firestore) getAllApplication(c *gin.Context) {
	iter := route.client.Collection("applications").Documents(route.ctx)
	ApplicationsData := []models.Application{}
	for {
		ApplicationData := models.Application{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		mapstructure.Decode(doc.Data(), &ApplicationData)
		ApplicationsData = append(ApplicationsData, ApplicationData)
	}
	c.JSON(http.StatusOK, ApplicationsData)
}

// Add new application to db
func (route *Firestore) createApplication(c *gin.Context) {
	var application models.Application
	c.BindJSON(&application)

	ref, _, err := route.client.Collection("applications").Add(route.ctx, application)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Success", "id": ref.ID})
}

// update existing application or create new if not existed
func (route *Firestore) updateApplicaiton(c *gin.Context) {
	id := c.Param("id")
	var application models.Application
	c.BindJSON(&application)
	_, err := route.client.Collection("applications").Doc(id).Set(route.ctx, application)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// Delele application from id given
func (route *Firestore) deleteApplication(c *gin.Context) {
	id := c.Param("id")
	_, err := route.client.Collection("applications").Doc(id).Delete(route.ctx)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Success"})
}
