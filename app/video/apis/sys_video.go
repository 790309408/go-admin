package apis

import (
	"fmt"
	"go-admin/app/video/models"
	"go-admin/app/video/service"
	"go-admin/app/video/service/dto"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

type SysVideo struct {
	api.Api
}

const (
	uploadDir      = "./uploads"
	maxFileSize    = 100 << 20 // 100MB
	defaultFPS     = 15
	defaultScale   = 480
	defaultSeconds = 10
)

/*
*
获取视频基本信息
首字母大写表示公开方法/导出
首字母小写表示私有方法/内部使用
*/
func (e SysVideo) GetVideoInfo(c *gin.Context) {

	// 接收上传文件
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传视频文件"})
		return
	}
	s := service.SysVideo{}
	video := &models.SysVideo{}
	s.GetVideoInfo(video, file, c)
	video.Format.Filename = file.Filename
	c.JSON(200, gin.H{
		"data": video,
	})
}

/*视频压缩*/
func (e SysVideo) CompressVideo(c *gin.Context) {
	startTime := time.Now()
	// 接收上传文件
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传视频文件"})
		return
	}
	qualityPost, _ := c.GetPostForm("quality")
	quality, _err := strconv.Atoi(qualityPost)
	if _err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入质量"})
		return
	}
	scale, _ := c.GetPostForm("scale")
	s := service.SysVideo{}
	reqParams := dto.CompressVideoReq{Quality: quality, Scale: scale}
	respData := s.CompressVideo(c, file, &reqParams)
	respData.TimeUsed = time.Since(startTime).String()
	c.JSON(200, gin.H{
		"data": respData,
	})

}

/*截取视频片段转换成gif*/
func (e SysVideo) CutVideoCoverToGif(c *gin.Context) {

	// 处理文件上传
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "视频文件不能为空"})
		return
	}
	reqParams := dto.ConvertRequest{}
	if err := c.ShouldBind(&reqParams); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 设置默认参数
	if reqParams.Duration <= 0 || reqParams.Duration > 60 {
		reqParams.Duration = defaultSeconds
	}
	if reqParams.FPS <= 0 || reqParams.FPS > 30 {
		reqParams.FPS = defaultFPS
	}
	if reqParams.Scale <= 0 || reqParams.Scale > 1920 {
		reqParams.Scale = defaultScale
	}
	s := service.SysVideo{}
	resp := s.CutVideoCoverToGif(c, file, &reqParams)
	c.JSON(http.StatusOK, gin.H{
		"download_url": resp.OutPutPath,
		"expires":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}

/*下载视频*/
func (e SysVideo) DownloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	fmt.Printf("filename:%s\n", filename)
	filePath := filepath.Join(uploadDir, filename)
	fmt.Printf("filePath:%s\n", filePath)
	// 验证文件存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}

	// 设置下载头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
	outputPath := filepath.Join(uploadDir, filename)
	defer os.Remove(outputPath) // 处理完成后删除压缩文件
}

// 示例
/*查询字符串参数*/
func (e SysVideo) GetQueryParams(c *gin.Context) {
	// 方式1：直接获取
	keyword := c.Query("param1")
	// 方式2：带默认值
	current := c.DefaultQuery("current", "1")
	size := c.DefaultQuery("size", "10")
	// 方式3：检查是否存在
	sort, exists := c.GetQuery("sort")
	if !exists {
		sort = "desc"
	}
	c.JSON(200, gin.H{
		"keyword": keyword,
		"current": current,
		"size":    size,
		"sort":    sort,
	})
}

/*表单参数*/
func (e SysVideo) GetFormParams(c *gin.Context) {
	// 方式1：直接获取
	username := c.PostForm("username")

	// 方式2：带默认值
	password := c.DefaultPostForm("password", "defaultPass")

	// 方式3：检查是否存在
	remember, exists := c.GetPostForm("remember")
	if !exists {
		remember = "false"
	}

	c.JSON(200, gin.H{
		"username": username,
		"password": password,
		"remember": remember,
	})
}

/*JSON 请求体*/
type User struct {
	// 表示必传，也不能传空
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func (e SysVideo) GetJsonParams(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "created",
		"data":   user,
	})
}

/*uri 参数绑定*/
type BookReq struct {
	// 表示必传，并且是uuid格式
	ID string `uri:"id" binding:"required,uuid"`
	// 必传，并且是字符串
	Type string `uri:"type" binding:"required"`
}

func (e SysVideo) GetUriParams(c *gin.Context) {
	var req BookReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"type": req.Type,
		"id":   req.ID,
	})
}

/*header 参数绑定*/
func (e SysVideo) GetHeaderParams(c *gin.Context) {
	// 获取 Authorization Header
	auth := c.GetHeader("Authorization")
	// 获取自定义 Header
	client := c.Request.Header.Get("X-Client-ID")
	c.JSON(200, gin.H{
		"auth":   auth,
		"client": client,
	})
}

/*文件上传*/
func (e SysVideo) UploadFiles(c *gin.Context) {
	// 单文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 保存文件
	dst := "./uploads/" + file.Filename
	c.SaveUploadedFile(file, dst)

	// 多文件上传
	form, _ := c.MultipartForm()
	files := form.File["files[]"]

	c.JSON(200, gin.H{
		"single_file": file.Filename,
		"file_count":  len(files),
	})
}

/*结构体参数绑定*/
type QueryParams struct {
	Page  int    `form:"page"` // 绑定查询参数
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
}

func (e SysVideo) GetStructParams(c *gin.Context) {
	var params QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"page": params.Page, "limit": params.Limit, "sort": params.Sort})
}

/*原始请求体*/
func (e SysVideo) GetRawBody(c *gin.Context) {
	// 读取原始数据
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.String(200, "Raw data: %s", string(data))
}

/*综合请求体*/
type ComplexRequest struct {
	ID      string `form:"id" uri:"id"`
	Query   string `form:"q"`
	Auth    string `header:"Authorization"`
	Content struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"content"`
}

func (e SysVideo) GetComprehensiveBinding(c *gin.Context) {
	var req ComplexRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error1": err.Error()})
		return
	}
	if err := c.ShouldBindHeader(&req); err != nil {
		c.JSON(400, gin.H{"error2": err.Error()})
		return
	}
	//根据 Content-Type 自动选择绑定方式
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error3": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": req})
}
