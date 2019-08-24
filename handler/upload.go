package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"netdisk/model"
	"netdisk/util"
	"os"
	"strconv"
)

//文件上传接口
func DoUploadHandler(c *gin.Context) {

	// 接收文件流及存储到本地目录
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		return
	}

	fileMeta := model.File{
		File_name: file.Filename,
		File_size: file.Size,
		File_addr: "D:\\tmp\\" + file.Filename,
		Status:    1,
	}

	err = c.SaveUploadedFile(file, fileMeta.File_addr)
	if err != nil {
		fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
		return
	}

	newFile, _ := os.Open(fileMeta.File_addr)
	if err != nil {
		fmt.Println("open error", err)
	}

	newFile.Seek(0, 0)
	fileMeta.File_sha1 = util.FileSha1(newFile)

	_ = fileMeta.OnFileUploadFinished()

	username := c.Request.FormValue("username")
	userFile := &model.UserFile{
		Username:  username,
		File_name: file.Filename,
		File_size: file.Size,
		File_sha1: fileMeta.File_sha1,
		Status:    1,
	}
	err = userFile.OnUserFileUpload()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "user: " + username + ", upload success",
			"code": 10001,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "user: " + username + ", upload failed",
			"code": 10000,
		})
	}

}

//文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")

	fmeta, _ := model.GetFileMeta(fsha1)

	file, err := os.Open(fmeta.File_addr)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{})
		fmt.Println("open da error", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{})
		fmt.Println("get data error", err)
		return
	}

	c.Header("content-disposition", "attachment; filename=\""+fmeta.File_name+"\"")
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// 查询批量的文件信息
func FileQueryHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")

	userFiles, err := model.QueryUserFileMetas(username, limit)
	if err != nil {
		log.Println(err.Error())
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		log.Println(err.Error())
	}
	c.Data(http.StatusOK, "application/json", data)
}

//秒传接口
func TryFastUploadHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _:= strconv.ParseInt(c.Request.FormValue("filesize"), 10, 64)

	fileMeta, err := model.GetFileMeta(filehash)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if fileMeta == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "秒传失败，请访问普通上传接口",
			"code": 10000,
		})
	}

	userFile := &model.UserFile{
		Username:  username,
		File_sha1: filehash,
		File_name: filename,
		File_size: filesize,
		Status: 1,
	}
	err = userFile.OnUserFileUpload()
	if err == nil {
		c.JSON(http.StatusOK,gin.H{
			"msg": "上传成功",
			"code": 10001,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "秒传失败,请稍后再试",
			"code": 10000,
		})
		return
	}
}
