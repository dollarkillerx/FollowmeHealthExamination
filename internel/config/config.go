package config

import (
	"encoding/json"
	"io/ioutil"
)

var CONFIG *Config

func InitConfig(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var conf Config
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return err
	}

	CONFIG = &conf
	return nil
}

type Config struct {
	Accounts    []Account   `json:"accounts"`
	Relation    []Relation  `json:"relation"`
	EmailServer EmailServer `json:"email_server"`
}

type Account struct {
	Uid   string `json:"uid"`
	Token string `json:"token"`
	Index string `json:"index"`
}

type Relation struct {
	Announcer string `json:"announcer"` // 发布者
	Follower  string `json:"follower"`  // 跟随者
	Email     string `json:"email"`     // 通知邮件
}

type EmailServer struct {
	SmtpServer string `json:"smtp_server"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
}
