package routers

import (
	"log/slog"

	_ "url-shortener/docs"
	"url-shortener/internal/handlers/url"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter initialize routing information
func InitRouter(log *slog.Logger, urlService *services.LinkService) *gin.Engine {
	r := gin.New()
	// Connect middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	apiv1 := r.Group("/api/v1")
	linksHandler := url.NewLinkHandler(log, urlService)

	link := apiv1.Group("/link")
	{
		link.POST("/", linksHandler.SaveLink)
		link.GET("/:link", linksHandler.GetLink)
	}

	return r
}
