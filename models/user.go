package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/xweis/go-gin/pkg/logging"
	"regexp"
)

type UserInfo struct {
	Model

	Uuid     string `gorm:"unique" json:"uuid"`
	Email    string `gorm:"unique" json:"name"`
	Password string `json:"password"`
	//1:开启， 0:关闭
	State int `json:"state"`
}

func AddUser(email string, password string, state int) error {
	//判断uuid 是否重复， 如果重复则重新生产, 3次
	for i := 0; i < 3; i++ {
		user := UserInfo{
			Uuid:     uuid.NewString(),
			Email:    email,
			Password: password,
			State:    state,
		}
		err := db.Create(&user).Error
		if err == nil {
			return nil
		}

		//判断uuid 是否重复， 如果重复则重新生产
		if !regexp.MustCompile(`Duplicate entry (.*) for key '(.*)_user_info.uuid'`).MatchString(err.Error()) {
			return err
		}
	}
	return errors.New("failed to create user after 3 attempts")
}

func GetUser(email string) (*UserInfo, error) {
	userInfo := UserInfo{}

	if err := db.Where("email = ?", email).First(&userInfo).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		logging.Error("%v", err)
		return &userInfo, err
	}

	return &userInfo, nil
}

func Login(email string, password string) (*UserInfo, error) {
	userInfo := UserInfo{}

	if err := db.Where("email = ? AND password = ?", email, password).First(&userInfo).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		logging.Error("%v", err)
		return &userInfo, err
	}

	logging.Info("email:%v, pass:%v, db: %v", email, password, userInfo)

	return &userInfo, nil
}
