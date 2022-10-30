package database

import "fmt"

// Img List
type User struct {
	ID       int64  `json:"id" gorm:"primary_key, column:id"`
	Token    string `json:"token" gorm:"varchar(128);index:idx_token_id;column:token"`
	Username string `json:"username" gorm:"varchar(64);column:username"`
	DateTime int64  `json:"datetime" gorm:"column:datetime"`
}

// TableName change table name
func (User) TableName() string {
	return "user"
}

func (user *User) Insert() (err error) {
	result := Eloquent.Create(&user)
	if result.Error != nil {
		err = result.Error
		fmt.Println(err.Error())
		return
	}
	return
}

// GetUser 列表
func (user *User) GetUser(token string) (users User, err error) {
	if err = Eloquent.First(&users, "token = ?", token).Error; err != nil {
		return
	}
	return
}

// FindName 列表
func (user *User) FindName(id int64) (users User, err error) {
	if err = Eloquent.First(&users, "id = ?", id).Error; err != nil {
		return
	}
	return
}
