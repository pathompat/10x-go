package main

import (
	"10x-go/models"
	"context"
	"log"
	"net/http"

	// "firebase/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type App struct {
	client *firestore.Client
	ctx    context.Context
}

func main() {
	route := App{}
	route.Init()

	router := gin.Default()

	// Get all applications coming in
	router.GET("/applications", func(c *gin.Context) {
		route.Application(c)
	})

	router.Run()

}

func (route *App) Init() {
	route.ctx = context.Background()
	sa := option.WithCredentialsFile("x-group-290609-firebase-adminsdk-wq05l-e29aa1fc60.json")
	app, err := firebase.NewApp(route.ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	route.client, err = app.Firestore(route.ctx)
	if err != nil {
		log.Fatalln(err)
	}

}

func (route *App) Application(c *gin.Context) {
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
