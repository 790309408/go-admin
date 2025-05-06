package dto

type CompressVideoReq struct {
	Quality int    `form:"quality"`
	Scale   string `form:"scale"`
}

type ConvertRequest struct {
	StartTime string `form:"start_time" binding:"required"` // 格式: 00:00:00
	Duration  int    `form:"duration"`                      // 秒数
	FPS       int    `form:"fps"`
	Scale     int    `form:"scale"`
}
