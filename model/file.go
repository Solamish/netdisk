package model

import "github.com/jinzhu/gorm"

type File struct {
	gorm.Model
	File_sha1 string
	File_name string
	File_size int64 `gorm:"default:'0'"`
	File_addr string
	Status    int `gorm:"default:'0'"`
}

func (file *File) OnFileUploadFinished() error {

	return DB.Create(file).Error
}

func GetFileMeta(fileSha1 string) (*File, error) {
	var file File
	err := DB.Where("file_sha1 = ?", fileSha1).First(&file).Error
	return &file, err
}