// Package router 注册 HTTP 路由。
package router

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"gin-enterprise/web"

	httpresponse "gin-enterprise/internal/handler/http/response"
	"gin-enterprise/internal/handler/http/websocket"
	"gin-enterprise/internal/service"
)

// OpenAPIJSONPath 是供前端代码生成工具读取的接口规范地址。
const OpenAPIJSONPath = "/openapi/json"

// Register 只挂载 HTTP 路由。业务依赖由 app 组装后传入。
func Register(app *gin.Engine, users *service.UserService, auth *service.AuthService, hub *websocket.Hub, categories *service.CategoryService, sites *service.SiteService) huma.API {
	httpresponse.RegisterErrors()

	config := huma.DefaultConfig("企业接口", "1.0.0")
	config.CreateHooks = nil
	config.Transformers = []huma.Transformer{httpresponse.FillRequest}
	config.Info.Description = "用户、认证、分类与站点接口。"
	api := humagin.New(app, config)
	app.GET(OpenAPIJSONPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, api.OpenAPI())
	})

	// /api 是全部业务接口的统一前缀。认证模块再挂 /auth，用户模块直接写 /users。
	group := huma.NewGroup(api, "/api")
	registerAuth(group, auth)
	registerUser(group, users)
	registerCategory(group, categories)
	registerSite(group, sites)

	websocket.Register(app, auth, hub)
	registerFrontend(app)

	return api
}

func registerFrontend(app *gin.Engine) {
	static, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		panic(err)
	}
	files := http.FileServer(http.FS(static))
	app.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}
		if _, err := fs.Stat(static, path.Clean(c.Request.URL.Path)); err == nil && c.Request.URL.Path != "/" {
			files.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.FileFromFS("/index.html", http.FS(static))
	})
}
