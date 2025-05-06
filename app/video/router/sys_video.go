package router

import (
	"go-admin/app/video/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

/*
Gin 框架中获取请求参数的 8 种常用方式
一、URL 路径参数
*/
func init() {
	routerCheckRole = append(routerCheckRole, registerSysVideoRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerSysVideoNoCheckRole)
}

// 需认证的路由代码
func registerSysVideoRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysVideo{}
	r := v1.Group("/sys-video").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{ // 获取视频
		r.POST("/get-video-info", api.GetVideoInfo)
		//视频压缩
		r.POST("/compress-video", api.CompressVideo)
		// 截取视频片段转换为gif图片
		r.POST("/cut-video-convert-gif", api.CutVideoCoverToGif)
	}
}

// 不需要认证的路由代码

func registerSysVideoNoCheckRole(v1 *gin.RouterGroup) {
	api := apis.SysVideo{}
	r := v1.Group("/sys-video")
	{ // 1.URL路径参数
		r.GET("/download/:filename", api.DownloadHandler)
		// 2.查询字符串参数
		r.GET("/getQuery", api.GetQueryParams)
		// 3.表单参数
		r.POST("/getFormParams", api.GetFormParams)
		// 4. JSON 请求体
		r.POST("/getJsonParams", api.GetJsonParams)
		// 5. URI 参数绑定
		r.GET("/getUriParams/:type/:id", api.GetUriParams)
		//6. Headers 参数
		r.GET("/getHeaderParams", api.GetHeaderParams)
		//7. 文件上传
		r.POST("/uploadFiles", api.UploadFiles)
		//8. 绑定结构体
		r.POST("/getStructParams", api.GetStructParams)
		//9. 原始请求体
		r.POST("/getRawBody", api.GetRawBody)
		//10. 综合绑定示例
		r.POST("/getComprehensiveBinding/:id", api.GetComprehensiveBinding)
	}
}
