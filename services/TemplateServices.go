package services

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/repositories"
)

func CreateTemplate(ctx *gin.Context, template models.MessageTemplate) (interface{}, error) {
	return repositories.CreateTemplate(ctx, template)
}
func GetTemplates(ctx *gin.Context, templateName string) (interface{}, error) {
	return repositories.GetTemplates(ctx, templateName)
}
func UpdateTemplate(ctx *gin.Context, templateId string, template models.MessageTemplate) (interface{}, error) {
	return repositories.UpdateTemplate(ctx, templateId, template)
}
func DeleteTemplate(ctx *gin.Context, templateId string) (interface{}, error) {
	return repositories.DeleteTemplate(ctx, templateId)
}

func GetAllTemplates(ctx *gin.Context) (interface{}, error) {
	return repositories.GetAllTemplates(ctx)
}
