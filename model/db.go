package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB



func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:zhy123@/file-server?charset=utf8&parseTime=True")
	if err == nil {

		DB = db
		db.AutoMigrate(&File{}, &User{},&Token{},&UserFile{})
		db.Model(&File{}).AddUniqueIndex("idx_file_hash","file_sha1")
		db.Model(&File{}).AddIndex("idx_status","status")
		db.Model(&User{}).AddUniqueIndex("idx_username","username")
		db.Model(&Token{}).AddUniqueIndex("idx_username","username")
		db.Model(&UserFile{}).AddUniqueIndex("idx_user_file","username","file_sha1")
		db.Model(&UserFile{}).AddIndex("idx_status","status")
		db.Model(&UserFile{}).AddIndex("idx_user_id","username")
		return db, err
	}
	return nil, err
}
