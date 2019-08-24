package handler

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	rPool "netdisk/cache/redis"
	"netdisk/model"
	"path"

	"os"
	"strconv"
	"strings"
	"time"
)

type MultiPartInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int //分块大小
	ChunkCount int //分块数量
}

//初始化分块上传
func InitialMultipartHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	//获取redis连接池
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	mpInfo := &MultiPartInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1025 * 1024))),
	}

	//将初始化信息存入缓存
	rConn.Do("HEST", "MP_"+mpInfo.UploadID, "chunkcount", mpInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "filesize", mpInfo.FileSize)
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "chunksize", mpInfo.ChunkSize)

	data, _ := json.Marshal(mpInfo)
	c.Data(http.StatusOK, "application/json", data)
}

//上传文件分块
func UploadPartHandler(c *gin.Context) {

	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": 0,
				"msg":  "Upload part failed",
				"data": nil,
			})
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//  更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
}

//通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	upid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//通过uploadid查询redis并判断所有分块是否上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "服务错误",
				"data": nil,
			})
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	// 不相等就是上传没有完成
	if totalCount != chunkCount {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "分块不完整",
				"data": nil,
			})
		return
	}

	//TODO 分块合并
	srcPath := "/data/fileserver_part/" + upid + "/"
	src, _ := os.Create(srcPath)
	destPath := "data/fileserver/"
	dst, _ := os.Create(destPath)
	io.Copy(dst, src)
	//更新唯一文件表和用户文件表
	fsize, _ := strconv.ParseInt((filesize), 10, 64)
	file := &model.File{
		File_size: fsize,
		File_name: filename,
		File_sha1: filehash,
		Status:    1,
		File_addr: destPath,
	}
	file.OnFileUploadFinished()
	userFile := &model.UserFile{
		Username:  username,
		File_sha1: filehash,
		File_name: filename,
		File_size: fsize,
		Status:    1,
	}
	userFile.OnUserFileUpload()

}
