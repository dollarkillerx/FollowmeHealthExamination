package server

import (
	"github.com/dollarkillerx/FollowmeHealthExamination/internel/config"
	"github.com/dollarkillerx/FollowmeHealthExamination/pkg/models"
	"github.com/dollarkillerx/FollowmeHealthExamination/utils"

	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Server struct {
	mu       sync.Mutex
	orderMap map[string][]models.Order
}

func NewServer() *Server {
	s := &Server{
		orderMap: map[string][]models.Order{},
	}

	go s.monitorFollowme()
	go s.orderVerification()

	return s
}

func (s *Server) monitorFollowme() {
	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ticker.C:
			for _, v := range config.CONFIG.Accounts {
				position, err := utils.CurrentPosition(v.Uid, v.Token)
				if err != nil {
					log.Println(err)
					continue
				}

				if position.Code != 0 {
					utils.Print(position)
					continue
				}

				go func() {
					s.mu.Lock()
					defer s.mu.Unlock()

					s.orderMap[v.Uid] = position.Data.Items
				}()

				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (s *Server) orderVerification() {
	ticker := time.NewTicker(time.Second * 4)
	for {
		select {
		case <-ticker.C:
			for _, v := range config.CONFIG.Relation {
				vv := v
				go s.orderVerificationInternal(vv)
			}
		}
	}
}

func (s *Server) getOrderByUid(uid string) ([]models.Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	orders, ex := s.orderMap[uid]
	return orders, ex
}

func sendEmail(to string, body string) error {
	format := time.Now().Format("2006-01-02 15:04:05")
	return utils.SendEmail(config.CONFIG.EmailServer.SmtpServer, config.CONFIG.EmailServer.Port, config.CONFIG.EmailServer.User, config.CONFIG.EmailServer.Password, to, "Followme Headlth Examination 预警: "+format, body)
}

func (s *Server) orderVerificationInternal(relation config.Relation) {
	announcer, ex := s.getOrderByUid(relation.Announcer)
	if !ex {
		err := sendEmail(relation.Email, fmt.Sprintf("信号源信息获取失败: %s", relation.Announcer))
		if err != nil {
			log.Println(err)
		}
		return
	}

	follower, ex := s.getOrderByUid(relation.Follower)
	if !ex {
		err := sendEmail(relation.Email, fmt.Sprintf("跟随者源信息获取失败: %s", relation.Announcer))
		if err != nil {
			log.Println(err)
		}
		return
	}

	var order []models.TwoWayOrder

	// 检查是否有 漏单
	for _, v := range announcer {
		exist := false
		for _, vv := range follower {
			if vv.Ex.TraderTradeID == v.TradeID {
				exist = true
			}
		}

		if !exist {
			tm := time.Unix(int64(v.OpenTime), 0)
			order = append(order, models.TwoWayOrder{
				AccountUid:     relation.Announcer,
				AccountOrderID: v.TradeID,
				Volume:         v.Volume,
				FollowmeUid:    relation.Follower,
				Status:         "漏单",
				CreateTime:     tm.Format("2006-01-02 15:04:05"),
			})
		}
	}

	// 检查是否有未平仓 订单
	for _, v := range follower {
		exist := false

		if v.Ex.TraderTradeID == 0 {
			continue
		}

		for _, vv := range announcer {
			if v.Ex.TraderTradeID == vv.TradeID {
				exist = true
			}
		}

		if !exist {
			tm := time.Unix(int64(v.OpenTime), 0)
			order = append(order, models.TwoWayOrder{
				AccountUid:     relation.Announcer,
				AccountOrderID: v.TradeID,
				Volume:         v.Volume,
				FollowmeUid:    relation.Follower,
				Status:         "未平订单",
				CreateTime:     tm.Format("2006-01-02 15:04:05"),
			})
		}
	}

	if len(order) == 0 {
		return
	}

	var newOder []models.TwoWayOrder

	// 限速器
	for _, v := range order {
		vv := v
		if !utils.GCache.Has(vv.AccountOrderID) {
			newOder = append(newOder, vv)
			utils.GCache.SetWithExpire(vv.AccountOrderID, "", time.Hour*2)
		}
	}

	if len(newOder) == 0 {
		return
	}

	// Send Email
	var tb = ""
	for _, v := range newOder {
		tb += fmt.Sprintf(table, v.AccountOrderID, v.Volume, v.CreateTime, v.FollowmeUid, v.AccountUid, v.Status)
	}

	em := strings.ReplaceAll(emailTemplate, "{thisIsBody}", tb)
	err := sendEmail(relation.Email, em)
	if err != nil {
		log.Println(err)
	}
}

var table = `<tr>
                                            <th>%s</th>
                                            <th>%.2f</th>
                                            <th>%s</th>
                                            <th>%s</th>
                                            <th>%s</th>
                                            <th>%s</th>
                                        </tr>`

var emailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Postman Followme Health Examination</title>

    <style>
        #outlook a {
            padding: 0
        }

        .ReadMsgBody {
            width: 100%
        }

        .ExternalClass {
            width: 100%
        }

        .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div {
            line-height: 100%
        }

        .mui-container-fixed {
            width: 600px;
            display: block;
            margin: 0 auto;
            clear: both;
            text-align: left;
            padding-left: 15px;
            padding-right: 15px
        }
    </style>

    <style>
        body {
            width: 100% !important;
            min-width: 100%;
            margin: 0;
            padding: 0
        }

        img {
            border: 0 none;
            height: auto;
            line-height: 100%;
            outline: 0;
            text-decoration: none
        }

        a img {
            border: 0 none
        }

        table {
            border-spacing: 0;
            border-collapse: collapse
        }

        td {
            padding: 0;
            text-align: left;
            word-break: break-word;
            -webkit-hyphens: auto;
            -moz-hyphens: auto;
            hyphens: auto;
            border-collapse: collapse !important
        }

        table, td {
            mso-table-lspace: 0;
            mso-table-rspace: 0
        }

        body, table, td, p, a, li, blockquote {
            -webkit-text-size-adjust: 100%;
            -ms-text-size-adjust: 100%
        }

        img {
            -ms-interpolation-mode: bicubic
        }

        body {
            color: #212121;
            font-family: "Helvetica Neue", Helvetica, Arial, Verdana, "Trebuchet MS";
            font-weight: 400;
            font-size: 14px;
            line-height: 1.429;
            letter-spacing: .001em;
            background-color: #FFF
        }

        a {
            color: #2196f3;
            text-decoration: none
        }

        p {
            margin: 0 0 10px
        }

        hr {
            color: #e0e0e0;
            background-color: #e0e0e0;
            height: 1px;
            border: 0
        }

        strong {
            font-weight: 700
        }

        h1, h2, h3 {
            margin-top: 20px;
            margin-bottom: 10px
        }

        h4, h5, h6 {
            margin-top: 10px;
            margin-bottom: 10px
        }

        .mui-body {
            margin: 0;
            padding: 0;
            height: 100%;
            width: 100%;
            color: #212121;
            font-family: "Helvetica Neue", Helvetica, Arial, Verdana, "Trebuchet MS";
            font-weight: 400;
            font-size: 14px;
            line-height: 1.429;
            letter-spacing: .001em;
            background-color: #FFF
        }

        .mui-btn {
            cursor: pointer;
            white-space: nowrap
        }

        a.mui-btn {
            font-weight: 500;
            font-size: 14px;
            color: #212121;
            line-height: 14px;
            letter-spacing: .03em;
            text-transform: uppercase;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent;
            border-top: 1px solid #FFF;
            border-left: 1px solid #FFF;
            border-right: 1px solid #FFF;
            border-bottom: 1px solid #FFF;
            color: #212121;
            background-color: #FFF;
            display: inline-block;
            text-decoration: none;
            text-align: center;
            border-radius: 3px;
            padding: 10px 25px;
            background-color: transparent
        }

        a.mui-btn.mui-btn--raised {
            border-top: 1px solid #f2f2f2;
            border-left: 1px solid #e6e6e6;
            border-right: 1px solid #e6e6e6;
            border-bottom: 2px solid #bababa
        }

        a.mui-btn.mui-btn--flat {
            background-color: transparent;
            color: #212121;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        a.mui-btn.mui-btn--primary {
            border-top: 1px solid #2196f3;
            border-left: 1px solid #2196f3;
            border-right: 1px solid #2196f3;
            border-bottom: 1px solid #2196f3;
            color: #FFF;
            background-color: #2196f3
        }

        a.mui-btn.mui-btn--primary.mui-btn--raised {
            border-top: 1px solid #51adf6;
            border-left: 1px solid #2196f3;
            border-right: 1px solid #2196f3;
            border-bottom: 2px solid #0a6ebd
        }

        a.mui-btn.mui-btn--primary.mui-btn--flat {
            background-color: transparent;
            color: #2196f3;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        a.mui-btn.mui-btn--danger {
            border-top: 1px solid #f44336;
            border-left: 1px solid #f44336;
            border-right: 1px solid #f44336;
            border-bottom: 1px solid #f44336;
            color: #FFF;
            background-color: #f44336
        }

        a.mui-btn.mui-btn--danger.mui-btn--raised {
            border-top: 1px solid #f77066;
            border-left: 1px solid #f44336;
            border-right: 1px solid #f44336;
            border-bottom: 2px solid #d2190b
        }

        a.mui-btn.mui-btn--danger.mui-btn--flat {
            background-color: transparent;
            color: #f44336;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        a.mui-btn.mui-btn--dark {
            border-top: 1px solid #424242;
            border-left: 1px solid #424242;
            border-right: 1px solid #424242;
            border-bottom: 1px solid #424242;
            color: #FFF;
            background-color: #424242
        }

        a.mui-btn.mui-btn--dark.mui-btn--raised {
            border-top: 1px solid #5c5c5c;
            border-left: 1px solid #424242;
            border-right: 1px solid #424242;
            border-bottom: 2px solid #1c1c1c
        }

        a.mui-btn.mui-btn--dark.mui-btn--flat {
            background-color: transparent;
            color: #424242;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        a.mui-btn.mui-btn--accent {
            border-top: 1px solid #ff4081;
            border-left: 1px solid #ff4081;
            border-right: 1px solid #ff4081;
            border-bottom: 1px solid #ff4081;
            color: #FFF;
            background-color: #ff4081
        }

        a.mui-btn.mui-btn--accent.mui-btn--raised {
            border-top: 1px solid #ff73a3;
            border-left: 1px solid #ff4081;
            border-right: 1px solid #ff4081;
            border-bottom: 2px solid #f30053
        }

        a.mui-btn.mui-btn--accent.mui-btn--flat {
            background-color: transparent;
            color: #ff4081;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        table.mui-btn > tr > td, table.mui-btn > tbody > tr > td {
            background-color: #FFF
        }

        table.mui-btn > tr > td > a, table.mui-btn > tbody > tr > td > a {
            color: #212121;
            border-top: 1px solid #FFF;
            border-left: 1px solid #FFF;
            border-right: 1px solid #FFF;
            border-bottom: 1px solid #FFF
        }

        table.mui-btn.mui-btn--raised > tr > td > a, table.mui-btn.mui-btn--raised > tbody > tr > td > a {
            border-top: 1px solid #f2f2f2;
            border-left: 1px solid #e6e6e6;
            border-right: 1px solid #e6e6e6;
            border-bottom: 2px solid #bababa
        }

        table.mui-btn.mui-btn--flat > tr > td, table.mui-btn.mui-btn--flat > tbody > tr > td {
            background-color: transparent
        }

        table.mui-btn.mui-btn--flat > tr > td > a, table.mui-btn.mui-btn--flat > tbody > tr > td > a {
            color: #212121;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        table.mui-btn > tr > td, table.mui-btn > tbody > tr > td {
            border-radius: 3px
        }

        table.mui-btn > tr > td > a, table.mui-btn > tbody > tr > td > a {
            font-weight: 500;
            font-size: 14px;
            color: #212121;
            line-height: 14px;
            letter-spacing: .03em;
            text-transform: uppercase;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent;
            display: inline-block;
            text-decoration: none;
            text-align: center;
            border-radius: 3px;
            padding: 10px 25px;
            background-color: transparent
        }

        table.mui-btn.mui-btn--primary > tr > td, table.mui-btn.mui-btn--primary > tbody > tr > td {
            background-color: #2196f3
        }

        table.mui-btn.mui-btn--primary > tr > td > a, table.mui-btn.mui-btn--primary > tbody > tr > td > a {
            color: #FFF;
            border-top: 1px solid #2196f3;
            border-left: 1px solid #2196f3;
            border-right: 1px solid #2196f3;
            border-bottom: 1px solid #2196f3
        }

        table.mui-btn.mui-btn--primary.mui-btn--raised > tr > td > a, table.mui-btn.mui-btn--primary.mui-btn--raised > tbody > tr > td > a {
            border-top: 1px solid #51adf6;
            border-left: 1px solid #2196f3;
            border-right: 1px solid #2196f3;
            border-bottom: 2px solid #0a6ebd
        }

        table.mui-btn.mui-btn--primary.mui-btn--flat > tr > td, table.mui-btn.mui-btn--primary.mui-btn--flat > tbody > tr > td {
            background-color: transparent
        }

        table.mui-btn.mui-btn--primary.mui-btn--flat > tr > td > a, table.mui-btn.mui-btn--primary.mui-btn--flat > tbody > tr > td > a {
            color: #2196f3;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        table.mui-btn.mui-btn--danger > tr > td, table.mui-btn.mui-btn--danger > tbody > tr > td {
            background-color: #f44336
        }

        table.mui-btn.mui-btn--danger > tr > td > a, table.mui-btn.mui-btn--danger > tbody > tr > td > a {
            color: #FFF;
            border-top: 1px solid #f44336;
            border-left: 1px solid #f44336;
            border-right: 1px solid #f44336;
            border-bottom: 1px solid #f44336
        }

        table.mui-btn.mui-btn--danger.mui-btn--raised > tr > td > a, table.mui-btn.mui-btn--danger.mui-btn--raised > tbody > tr > td > a {
            border-top: 1px solid #f77066;
            border-left: 1px solid #f44336;
            border-right: 1px solid #f44336;
            border-bottom: 2px solid #d2190b
        }

        table.mui-btn.mui-btn--danger.mui-btn--flat > tr > td, table.mui-btn.mui-btn--danger.mui-btn--flat > tbody > tr > td {
            background-color: transparent
        }

        table.mui-btn.mui-btn--danger.mui-btn--flat > tr > td > a, table.mui-btn.mui-btn--danger.mui-btn--flat > tbody > tr > td > a {
            color: #f44336;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        table.mui-btn.mui-btn--dark > tr > td, table.mui-btn.mui-btn--dark > tbody > tr > td {
            background-color: #424242
        }

        table.mui-btn.mui-btn--dark > tr > td > a, table.mui-btn.mui-btn--dark > tbody > tr > td > a {
            color: #FFF;
            border-top: 1px solid #424242;
            border-left: 1px solid #424242;
            border-right: 1px solid #424242;
            border-bottom: 1px solid #424242
        }

        table.mui-btn.mui-btn--dark.mui-btn--raised > tr > td > a, table.mui-btn.mui-btn--dark.mui-btn--raised > tbody > tr > td > a {
            border-top: 1px solid #5c5c5c;
            border-left: 1px solid #424242;
            border-right: 1px solid #424242;
            border-bottom: 2px solid #1c1c1c
        }

        table.mui-btn.mui-btn--dark.mui-btn--flat > tr > td, table.mui-btn.mui-btn--dark.mui-btn--flat > tbody > tr > td {
            background-color: transparent
        }

        table.mui-btn.mui-btn--dark.mui-btn--flat > tr > td > a, table.mui-btn.mui-btn--dark.mui-btn--flat > tbody > tr > td > a {
            color: #424242;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        table.mui-btn.mui-btn--accent > tr > td, table.mui-btn.mui-btn--accent > tbody > tr > td {
            background-color: #ff4081
        }

        table.mui-btn.mui-btn--accent > tr > td > a, table.mui-btn.mui-btn--accent > tbody > tr > td > a {
            color: #FFF;
            border-top: 1px solid #ff4081;
            border-left: 1px solid #ff4081;
            border-right: 1px solid #ff4081;
            border-bottom: 1px solid #ff4081
        }

        table.mui-btn.mui-btn--accent.mui-btn--raised > tr > td > a, table.mui-btn.mui-btn--accent.mui-btn--raised > tbody > tr > td > a {
            border-top: 1px solid #ff73a3;
            border-left: 1px solid #ff4081;
            border-right: 1px solid #ff4081;
            border-bottom: 2px solid #f30053
        }

        table.mui-btn.mui-btn--accent.mui-btn--flat > tr > td, table.mui-btn.mui-btn--accent.mui-btn--flat > tbody > tr > td {
            background-color: transparent
        }

        table.mui-btn.mui-btn--accent.mui-btn--flat > tr > td > a, table.mui-btn.mui-btn--accent.mui-btn--flat > tbody > tr > td > a {
            color: #ff4081;
            border-top: 1px solid transparent;
            border-left: 1px solid transparent;
            border-right: 1px solid transparent;
            border-bottom: 1px solid transparent
        }

        a.mui-btn--small, table.mui-btn--small > tr > td > a, table.mui-btn--small > tbody > tr > td > a {
            font-size: 13px;
            padding: 7.8px 15px
        }

        a.mui-btn--large, table.mui-btn--large > tr > td > a, table.mui-btn--large > tbody > tr > td > a {
            font-size: 14px;
            padding: 19px 25px
        }

        .mui-container, .mui-container-fixed {
            max-width: 600px;
            display: block;
            margin: 0 auto;
            clear: both;
            text-align: left;
            padding-left: 15px;
            padding-right: 15px
        }

        .mui-container-fixed {
            width: 600px
        }

        .mui-divider {
            display: block;
            height: 1px;
            background-color: #e0e0e0
        }

        .mui--divider-top {
            border-top: 1px solid #e0e0e0
        }

        .mui--divider-bottom {
            border-bottom: 1px solid #e0e0e0
        }

        .mui--divider-left {
            border-left: 1px solid #e0e0e0
        }

        .mui--divider-right {
            border-right: 1px solid #e0e0e0
        }

        .mui-panel {
            padding: 15px;
            border-radius: 0;
            background-color: #FFF;
            border-top: 1px solid #ededed;
            border-left: 1px solid #e6e6e6;
            border-right: 1px solid #e6e6e6;
            border-bottom: 2px solid #d4d4d4
        }

        .mui--text-left {
            text-align: left
        }

        .mui--text-right {
            text-align: right
        }

        .mui--text-center {
            text-align: center
        }

        .mui--text-justify {
            text-align: justify
        }

        .mui-image--fix {
            display: block
        }

        .mui--text-dark {
            color: #212121
        }

        .mui--text-dark-secondary {
            color: #757575
        }

        .mui--text-dark-hint {
            color: #9e9e9e
        }

        .mui--text-light {
            color: #FFF
        }

        .mui--text-light-secondary {
            color: #b3b3b3
        }

        .mui--text-light-hint {
            color: gray
        }

        .mui--text-accent {
            color: #ff4081
        }

        .mui--text-accent-secondary {
            color: #ff82ad
        }

        .mui--text-accent-hint {
            color: #ffa6c4
        }

        .mui--text-black {
            color: #000
        }

        .mui--text-white {
            color: #FFF
        }

        .mui--text-danger {
            color: #f44336
        }

        .mui--text-display4 {
            font-weight: 300;
            font-size: 112px;
            line-height: 112px
        }

        .mui--text-display3 {
            font-weight: 400;
            font-size: 56px;
            line-height: 56px
        }

        .mui--text-display2 {
            font-weight: 400;
            font-size: 45px;
            line-height: 48px
        }

        .mui--text-display1, h1 {
            font-weight: 400;
            font-size: 34px;
            line-height: 40px
        }

        .mui--text-headline, h2 {
            font-weight: 400;
            font-size: 24px;
            line-height: 32px
        }

        .mui--text-title, h3 {
            font-weight: 400;
            font-size: 20px;
            line-height: 28px
        }

        .mui--text-subhead, h4 {
            font-weight: 400;
            font-size: 16px;
            line-height: 24px
        }

        .mui--text-body2, h5 {
            font-weight: 500;
            font-size: 14px;
            line-height: 24px
        }

        .mui--text-body1 {
            font-weight: 400;
            font-size: 14px;
            line-height: 20px
        }

        .mui--text-caption {
            font-weight: 400;
            font-size: 12px;
            line-height: 16px
        }

        .mui--text-menu {
            font-weight: 500;
            font-size: 13px;
            line-height: 17px
        }

        .mui--text-button {
            font-weight: 500;
            font-size: 14px;
            line-height: 18px;
            text-transform: uppercase
        }
    </style>

    <style>
        #content-wrapper h2 {
            margin-top: 0px;
            margin-bottom: 0px;
        }

        #content-wrapper > tbody > tr > td {
            padding-bottom: 15px;
        }

        #content-wrapper .mui--divider-top {
            padding-top: 15px;
        }

        #last-cell {
            padding-bottom: 15px;
        }
    </style>
</head>
<body>
<table class="mui-body" cellpadding="0" cellspacing="0" border="0">
    <tr>
        <td>
            <center>
                <!--[if mso]>
                <table>
                    <tr>
                        <td class="mui-container-fixed"><![endif]-->
                <div class="mui-container">
                    <!--

                    email goes here

                    -->

                    <h3 class="mui--text-center">Postman Followme Health Examination</h3>
                    <table cellpadding="0" cellspacing="0" border="0" width="100%">
                        <tr>
                            <td class="mui-panel">
                                <table
                                        id="content-wrapper"
                                        border="0"
                                        cellpadding="0"
                                        cellspacing="0"
                                        width="100%"
                                >
                                    <tbody>
                                    <tr>
                                        <td>
                                            <h2>Postman Followme Health Examination 订单状态预警!</h2>
                                        </td>
                                    </tr>

                                    <table class="mui-table">
                                        <thead>
                                        <tr>
                                            <th>订单ID</th>
                                            <th>手数</th>
                                            <th>开仓时间</th>
                                            <th>跟随者ID</th>
                                            <th>信号源ID</th>
                                            <th>状态</th>
                                        </tr>
                                        </thead>
                                        <tbody>
                                        {thisIsBody}
                                        </tbody>
                                    </table>

                                </table>
                            </td>
                        </tr>
                    </table>
                </div>
                <!--[if mso]></td></tr></table><![endif]-->
            </center>
        </td>
    </tr>
</table>
</body>
</html>
`
