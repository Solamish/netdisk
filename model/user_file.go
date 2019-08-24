package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserFile struct {
	gorm.Model
	Username  string
	File_sha1 string
	File_name string
	File_size int64 `gorm:"default:'0'"`
	Status    int   `gorm:"default:'0'"`
}

func (userfile *UserFile) OnUserFileUpload() error {
	return DB.Create(userfile).Error
}

// 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	rows, err := DB.Exec("select file_sha1, file_name, file_size, created_at," +
		" updated_at from user_files where username = ? and limit = ?",username, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.File_sha1, &ufile.File_name, &ufile.File_size,
			&ufile.CreatedAt, &ufile.UpdatedAt)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
