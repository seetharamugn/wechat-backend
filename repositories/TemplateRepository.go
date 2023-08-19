package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"time"
)

var templateCollection *mongo.Collection = initializers.OpenCollection(initializers.Client, "message_templates")

func CreateTemplate(ctx *gin.Context, template models.MessageTemplate) (interface{}, error) {
	templateId := GenerateRandomTemplateId()
	newTemplate := models.MessageTemplate{
		TemplateId: templateId,
		Name:       template.Name,
		Category:   template.Category,
		Content:    template.Content,
		Status:     "INREVIEW",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	var existingTemplate models.MessageTemplate
	err := templateCollection.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"Name": template.Name},
	}}).Decode(&existingTemplate)
	if err == nil {
		// Either username or email already exists
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "templateName already exists",
		})
		ctx.Abort()
		return nil, err
	} else if err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal server error",
		})
		ctx.Abort()
		return nil, err
	}
	resp, err := templateCollection.InsertOne(context.TODO(), newTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create template",
		})
		ctx.Abort()
		return nil, err
	}
	return resp, nil
}
func GenerateRandomTemplateId() int {
	randNumber := 10000000 + rand.Intn(99999999-10000000)
	var template models.MessageTemplate
	err := templateCollection.FindOne(context.TODO(), bson.M{"templateId": randNumber}).Decode(&template)
	if err != nil {
		return randNumber
	}
	return GenerateRandomTemplateId()
}

func GetTemplates(ctx *gin.Context, name string) (interface{}, error) {
	var templates models.MessageTemplate
	if name != "" {
		err := templateCollection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&templates)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, Dao.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get templates",
				Data:       err.Error(),
			})
			ctx.Abort()
			return templates, err
		}
	}
	return templates, nil
}
func GetAllTemplates(ctx *gin.Context) (interface{}, error) {
	var templates []models.MessageTemplate
	cursor, err := templateCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get templates",
			Data:       err.Error(),
		})
		ctx.Abort()
		return templates, err
	}
	for cursor.Next(context.TODO()) {
		var template models.MessageTemplate
		err = cursor.Decode(&template)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func UpdateTemplate(ctx *gin.Context, templateId string, template models.MessageTemplate) (interface{}, error) {
	var existingTemplate models.MessageTemplate
	err := templateCollection.FindOne(context.TODO(), bson.M{"templateId": templateId}).Decode(&existingTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get templates",
			Data:       err.Error(),
		})
		ctx.Abort()
		return nil, err
	}
	newTemplate := models.MessageTemplate{
		Name:      template.Name,
		Category:  template.Category,
		Content:   template.Content,
		Status:    template.Status,
		Language:  template.Language,
		Footer:    template.Footer,
		UpdatedAt: time.Now(),
	}
	resp, err := templateCollection.ReplaceOne(context.TODO(), bson.M{"templateId": templateId}, newTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update template",
			Data:       err.Error(),
		})
		ctx.Abort()
		return nil, err
	}
	return resp, nil
}
func DeleteTemplate(ctx *gin.Context, templateId string) (interface{}, error) {
	var existingTemplate models.MessageTemplate
	err := templateCollection.FindOne(context.TODO(), bson.M{"templateId": templateId}).Decode(&existingTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get templates",
			Data:       err.Error(),
		})
		return nil, err
	}
	resp, err := templateCollection.DeleteOne(context.TODO(), bson.M{"templateId": templateId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Dao.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete template",
			Data:       err.Error(),
		})
		return nil, err
	}
	return resp, nil
}
