package controllers

import (
	"context"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-playground/validator/v10"

	"github.com/milaabl/proof-of-sloth-api/db"
	"github.com/milaabl/proof-of-sloth-api/models"
	"github.com/milaabl/proof-of-sloth-api/responses"
)

var albumCollection *mongo.Collection = db.GetCollection(db.DB, "albums")
var validate = validator.New()

func CreateAlbum() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var album models.Album
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&album); err != nil {
            c.IndentedJSON(http.StatusBadRequest, responses.AlbumDTO{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&album); validationErr != nil {
            c.IndentedJSON(http.StatusBadRequest, responses.AlbumDTO{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        newAlbum := models.Album{
            Id: primitive.NewObjectID(),
            Title: album.Title,
            Price: album.Price,
            Artist: album.Artist,
        }

        result, err := albumCollection.InsertOne(ctx, newAlbum)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.IndentedJSON(http.StatusCreated, responses.AlbumDTO{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}

func GetAlbumByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        albumId := c.Param("id")
        var album models.Album
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(albumId)

        err := albumCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&album)
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.IndentedJSON(http.StatusOK, responses.AlbumDTO{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": album}})
    }
}

func DeleteAlbum() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        albumId := c.Param("id")
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(albumId)

        result, err := albumCollection.DeleteOne(ctx, bson.M{"id": objId})
        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.IndentedJSON(http.StatusOK, responses.AlbumDTO{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}

func EditAlbum() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        albumId := c.Param("id")
        var album models.Album
        defer cancel()
        objId, _ := primitive.ObjectIDFromHex(albumId)

        //validate the request body
        if err := c.BindJSON(&album); err != nil {
            c.JSON(http.StatusBadRequest, responses.AlbumDTO{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&album); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.AlbumDTO{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        update := bson.M{"title": album.Title, "artist": album.Artist, "price": album.Price}
        result, err := albumCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated user details
        var updatedUser models.Album
        if result.MatchedCount == 1 {
            err := albumCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
            if err != nil {
                c.JSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.AlbumDTO{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
    }
}

func GetAlbums() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

        var albums []models.Album
        defer cancel()

        results, err := albumCollection.Find(ctx, bson.M{})

        if err != nil {
            c.IndentedJSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

		defer results.Close(ctx)
        for results.Next(ctx) {
            var singleAlbum models.Album
            if err = results.Decode(&singleAlbum); err != nil {
                c.JSON(http.StatusInternalServerError, responses.AlbumDTO{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }

            albums = append(albums, singleAlbum)
        }

        c.JSON(http.StatusOK,
            responses.AlbumDTO{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": albums}},
        )
	}	
}