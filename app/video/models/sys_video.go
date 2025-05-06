package models

// 视频元数据结构
type SysVideo struct {
	Format struct {
		Filename   string `json:"filename"`
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"` // 单位：秒
		Size       string `json:"size"`     // 文件大小
		BitRate    string `json:"bit_rate"` // 比特率
	} `json:"format"`
	VideoStream struct {
		CodecName string `json:"codec_name"` // 编码格式
		Width     int    `json:"width"`      // 分辨率宽
		Height    int    `json:"height"`     // 分辨率高
		FrameRate string `json:"frame_rate"` // 帧率
		BitRate   string `json:"bit_rate"`   // 视频比特率
	} `json:"video_stream"`
	AudioStream struct {
		CodecName  string `json:"codec_name"`  // 音频编码
		SampleRate string `json:"sample_rate"` // 采样率
		Channels   int    `json:"channels"`    // 声道数
		BitRate    string `json:"bit_rate"`    // 音频比特率
	} `json:"audio_stream"`
}
