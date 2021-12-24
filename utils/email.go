package utils

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

func SendEmail(host string, port int, username, password, to, title, body string) error {
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", "Postman"+"<"+username+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", to)                           //发送给多个用户
	m.SetHeader("Subject", title)                   //设置邮件主题
	m.SetBody("text/html", body)                    //设置邮件正文

	err := d.DialAndSend(m)
	return err
}
