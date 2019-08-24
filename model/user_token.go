package model

type Token struct {
	ID	int `gorm:"primary_key;AUTO_INCREMENT"`
	Username string
	User_token string
}

func UpdateToken(username string, user_token string) error {
	return DB.Exec("replace into tokens (`username`,`user_token`) values (?,?)",username,user_token).Error
}