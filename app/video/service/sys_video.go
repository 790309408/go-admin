package service

import (
	"encoding/json"
	"fmt"
	"go-admin/app/video/models"
	"go-admin/app/video/models/vo"
	"go-admin/app/video/service/dto"
	"go-admin/common/service"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uploadDir = "./uploads"
)

type SysVideo struct {
	service.Service
	FFmpegPath  string // FFmpeg可执行文件路径
	FFprobePath string // FFprobe可执行文件路径
}
type Result struct {
	Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
		Size       string `json:"size"`
		BitRate    string `json:"bit_rate"`
	} `json:"format"`
	Streams []struct {
		CodecType  string `json:"codec_type"`
		CodecName  string `json:"codec_name"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		BitRate    string `json:"bit_rate"`
		RFrameRate string `json:"r_frame_rate"`
		SampleRate string `json:"sample_rate"`
		Channels   int    `json:"channels"`
	} `json:"streams"`
}

// 创建新的视频处理器
func newVideoProcessor() *SysVideo {
	//获取ffmpeg路径
	exePath, err := getExecPath()
	if err != nil {
		fmt.Println("获取ffmpeg路径错误:", err)
	}
	return &SysVideo{
		FFmpegPath:  exePath + "/ffmpeg/bin/ffmpeg.exe",
		FFprobePath: exePath + "/ffmpeg/bin/ffprobe.exe",
	}
}

/*清除文件*/
func cleanupFiles(files ...string) {
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			os.Remove(f)
		}
	}
}

/*解析视频文件信息*/
func (vp *SysVideo) parseVideoInfo(inputPath string) (resultData Result, err error) {
	result := Result{}
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return result, fmt.Errorf("输入文件不存在: %s", inputPath)
	}
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}
	cmd := exec.Command(vp.FFprobePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("获取视频信息失败: %v", err)
	}
	// 解析 JSON 输出
	fmt.Println("output", string(output))
	if err := json.Unmarshal(output, &result); err != nil {
		return result, fmt.Errorf("解析视频信息失败: %v", err)
	}
	return result, nil
}

/*压缩视频*/
func (vp *SysVideo) compressVideo(inputPath string, outputPath string, req *dto.CompressVideoReq) (compressData vo.CompressVideoResponse, err error) {
	result := vo.CompressVideoResponse{}
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return result, fmt.Errorf("输入文件不存在: %s", inputPath)
	}
	args := []string{"-i", inputPath, "-c:v", "libx264"}
	// 添加质量参数
	if req.Quality > 0 {
		args = append(args, "-crf", strconv.Itoa(req.Quality))
	}
	// 添加缩放参数
	if req.Scale != "" {
		args = append(args, "-vf", "scale="+req.Scale)
	}
	// 添加输出文件路径
	args = append(args, "-y", outputPath)
	cmd := exec.Command(vp.FFmpegPath, args...)
	// output, err := cmd.CombinedOutput()
	// 执行压缩
	if err := cmd.Run(); err != nil {
		return result, fmt.Errorf("压缩视频失败: %v", err)

	}
	fmt.Printf("视频压缩完成: %s\n", outputPath)
	// 获取文件信息
	origInfo, _ := os.Stat(inputPath)
	compressedInfo, _ := os.Stat(outputPath)
	result.OriginSize = strconv.Itoa(int(origInfo.Size()))
	result.CompressSize = strconv.Itoa(int(compressedInfo.Size()))
	return result, nil

}

/*截取视频转换为gif图*/
func (vp *SysVideo) cutVideoConvertGif(inputPath, outputPath string, req *dto.ConvertRequest) (err error) {
	palettePath := filepath.Join(uploadDir, "palette.png")
	defer os.Remove(palettePath)
	fmt.Println("palettePath:", palettePath)
	fmt.Println("outputPath:", outputPath)
	fmt.Println("inputPath:", inputPath)
	// 生成调色板
	paletteCmd := exec.Command(vp.FFmpegPath,
		"-ss", req.StartTime,
		"-t", fmt.Sprintf("%d", req.Duration),
		"-i", inputPath,
		"-vf", fmt.Sprintf("fps=%d,scale=%d:-1:flags=lanczos,palettegen", req.FPS, req.Scale),
		"-y", palettePath,
	)
	if output, err := paletteCmd.CombinedOutput(); err != nil {
		fmt.Println("paletteCmd output1:", string(output))
		return fmt.Errorf("调色板生成失败: %s", string(output))
	}

	// 生成最终GIF
	gifCmd := exec.Command(vp.FFmpegPath,
		"-ss", req.StartTime,
		"-t", fmt.Sprintf("%d", req.Duration),
		"-i", inputPath,
		"-i", palettePath,
		"-filter_complex", fmt.Sprintf("fps=%d,scale=%d:-1:flags=lanczos[x];[x][1:v]paletteuse", req.FPS, req.Scale),
		"-y", outputPath,
	)
	if output, err := gifCmd.CombinedOutput(); err != nil {
		fmt.Println("paletteCmd output2:", string(output))
		return fmt.Errorf("GIF生成失败: %s", string(output))
	}

	return nil

}

/*获取执行目录路径*/
func getExecPath() (string, error) {

	exePath, err := os.Getwd()
	log.Println("exePath:", exePath)
	if err != nil {
		return "", err
	}
	return exePath, nil
}

/*获取视频信息*/
func (e *SysVideo) GetVideoInfo(video *models.SysVideo, file *multipart.FileHeader, c *gin.Context) {
	if video == nil {
		log.Println("video is nil")
		return
	}
	// 创建临时文件
	os.MkdirAll(uploadDir, 0755)
	tempPath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))

	// 保存临时文件
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}
	defer os.Remove(tempPath) // 处理完成后删除临时文件
	processor := newVideoProcessor()
	parseData, err := processor.parseVideoInfo(tempPath)
	if err != nil {
		log.Println("获取视频信息失败:", err)
		return
	}
	video.Format.Duration = parseData.Format.Duration
	video.Format.FormatName = parseData.Format.FormatName
	video.Format.Size = parseData.Format.Size
	video.Format.BitRate = parseData.Format.BitRate
	//视频信息
	video.VideoStream.CodecName = parseData.Streams[0].CodecName
	video.VideoStream.Width = parseData.Streams[0].Width
	video.VideoStream.Height = parseData.Streams[0].Height
	video.VideoStream.BitRate = parseData.Streams[0].BitRate
	video.VideoStream.FrameRate = parseData.Streams[0].RFrameRate
	//音频信息
	video.AudioStream.CodecName = parseData.Streams[1].CodecName
	video.AudioStream.Channels = parseData.Streams[1].Channels
	video.AudioStream.SampleRate = parseData.Streams[1].SampleRate
	video.AudioStream.BitRate = parseData.Streams[1].BitRate

}

/*
**视频压缩
返回值：值类型还是指针类型
小型结构（<100B）：值类型性能更优
大型结构（>1KB）：指针类型性能优势明显
**
*/
func (e *SysVideo) CompressVideo(c *gin.Context, file *multipart.FileHeader, req *dto.CompressVideoReq) (compressData *vo.CompressVideoResponse) {
	fmt.Println("reqParams:", req)
	if req == nil {
		log.Println("req is nil")
		return
	}
	// 创建临时文件
	os.MkdirAll(uploadDir, 0755)
	tempPath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
	// 保存临时文件
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}
	defer os.Remove(tempPath) // 处理完成后删除临时文件
	outputFilename := "compressed_" + file.Filename
	outputPath := filepath.Join(uploadDir, outputFilename)
	// defer os.Remove(outputPath) // 处理完成后删除压缩文件
	// 创建视频处理器
	processor := newVideoProcessor()
	resp, _ := processor.compressVideo(tempPath, outputPath, req)
	resp.OutPutPath = "/download/" + outputFilename
	return &resp
}

/**
**视频裁剪转成gif图片
 */
func (e *SysVideo) CutVideoCoverToGif(c *gin.Context, file *multipart.FileHeader, req *dto.ConvertRequest) (cutData vo.CutVideoCoverToGifResponse) {

	if req == nil {
		log.Println("req is nil")
		return
	}
	// 生成唯一ID防止冲突
	taskID := uuid.New().String()
	tempVideo := filepath.Join(uploadDir, taskID+".mp4")
	defer os.Remove(tempVideo)
	// 保存临时文件
	if err := c.SaveUploadedFile(file, tempVideo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}
	defer os.Remove(tempVideo) // 处理完成后删除临时文件
	outputFile := filepath.Join(uploadDir, fmt.Sprintf("%s.gif", taskID))
	// 生成调色板（优化颜色）
	paletteFile := filepath.Join(uploadDir, taskID+"_palette.png")
	defer os.Remove(paletteFile)
	// 创建视频处理器
	processor := newVideoProcessor()
	processor.cutVideoConvertGif(tempVideo, outputFile, req)
	respData := vo.CutVideoCoverToGifResponse{}

	respData.OutPutPath = fmt.Sprintf("/download/%s.gif", taskID)
	return respData
}
