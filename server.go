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
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	f := Firestore{}
	f.Init()

	router := gin.Default()
	router.Use(gin.BasicAuth(gin.Accounts{
		"admin": password,
	}))

	applicationRoute := router.Group("/applications")
	{
		applicationRoute.GET("/", f.getAllApplication)
	}

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
