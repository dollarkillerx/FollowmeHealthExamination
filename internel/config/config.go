package config

type Config struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Uid   string `json:"uid"`
	Token string `json:"token"`
}

type Relation struct {
	Announcer string `json:"announcer"` // 发布者
	Follower  string `json:"follower"`  // 跟随者
	Email     string `json:"email"`     // 通知邮件
}

type EmailServer struct {
	SmtpServer string `json:"smtp_server"`
	Port       int    `json:"port"`
}
