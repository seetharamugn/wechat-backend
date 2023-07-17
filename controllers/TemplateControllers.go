package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/seetharamugn/wachat/Dao"
	"github.com/seetharamugn/wachat/models"
	"github.com/seetharamugn/wachat/services"
	"net/http"
)

func CreateTemplate(ctx *gin.Context) {
	var body models.MessageTemplate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	if body.TemplateName == "" || body.TemplateCategory == "" || body.TemplateContent == "" || body.TemplateLanguage == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Template"})
		ctx.Abort()
		return
	}
	resp, _ := services.CreateTemplate(ctx, body)
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "Template created successfully",
		Data:       resp,
	})
}
func GetTemplate(ctx *gin.Context) {
	templateName := ctx.Param("name")
	resp, _ := services.GetTemplates(ctx, templateName)
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       resp,
	})
}

func GetTemplateList(ctx *gin.Context) {
	resp, _ := services.GetAllTemplates(ctx)
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       resp,
	})
}

func UpdateTemplate(ctx *gin.Context) {
	templateId := ctx.Param("templateId")
	var body models.MessageTemplate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, Dao.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request body",
			Data:       err.Error(),
		})
		ctx.Abort()
		return
	}
	resp, _ := services.UpdateTemplate(ctx, templateId, body)
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       resp,
	})
}

func DeleteTemplate(ctx *gin.Context) {
	templateId := ctx.Param("templateId")
	resp, _ := services.DeleteTemplate(ctx, templateId)
	ctx.JSON(http.StatusOK, Dao.Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       resp,
	})
}
