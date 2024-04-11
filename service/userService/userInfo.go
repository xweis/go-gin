package userService

import (
	"errors"
	"github.com/xweis/go-gin/models"
	"github.com/xweis/go-gin/pkg/logging"
	"github.com/xweis/go-gin/pkg/util/aes"
	"github.com/xweis/go-gin/pkg/util/jwt"
)

type User struct {
	Email    string
	Password string
	Id       string
}

func (user *User) AddUser() error {
	encodePassword, err := aes.EncryptByAes([]byte(user.Password))
	if err != nil {
		logging.Error("%v", err)
		return err
	}

	logging.Info("%s", encodePassword)
	if err := models.AddUser(user.Email, encodePassword, 1); err != nil {
		logging.Error("%v", err)
		return err
	}

	return nil
}

func (user *User) GetUser() (*models.UserInfo, error) {

	userInfo, err := models.GetUser(user.Email)
	if err != nil {
		logging.Error("%v", err)
		return userInfo, err
	}

	return userInfo, nil
}

func (user *User) Login() (jwt.TokenInfo, error) {
	token := jwt.TokenInfo{}
	encodePassword, err := aes.EncryptByAes([]byte(user.Password))
	if err != nil {
		logging.Error("%v", err)
		return token, err
	}

	//数据库查询账号密码
	userInfo, err := models.Login(user.Email, encodePassword)
	if err != nil {
		logging.Error("%v", err)
		return token, err
	}

	if userInfo.Email == user.Email {
		token, err = jwt.GetJwtToken(userInfo.Email)
		if err != nil {
			return jwt.TokenInfo{}, nil
		}

		return token, nil
	}

	return token, errors.New("账号验证失败")
}
