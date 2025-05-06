package test

import (
	"testing"

	"go-admin/common/video"
)

// 执行测试命令 go test video_test.go
func TestCompressVideo(t *testing.T) {
	vp := video.NewVideoProcessor("../ffmpeg/bin/ffmpeg.exe")
	err := vp.CompressVideo("../static/video/example.mp4", "../static/video/test_compressed.mp4", 20, "")
	if err != nil {
		t.Errorf("压缩视频失败: %v", err)
	}
}
