package vo

// 压缩响应数据
type CompressVideoResponse struct {
	OriginSize   string `json:"origin_size"`
	CompressSize string `json:"compress_size"`
	OutPutPath   string `json:"out_put_path"`
	TimeUsed     string `json:"time_used"`
}

type CutVideoCoverToGifResponse struct {
	OutPutPath string `json:"out_put_path"`
}
