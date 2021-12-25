package server

import (
	"encoding/json"
	"fmt"
	"github.com/dollarkillerx/FollowmeHealthExamination/internel/config"
	"github.com/dollarkillerx/FollowmeHealthExamination/pkg/models"
	"github.com/dollarkillerx/FollowmeHealthExamination/utils"
	"log"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ticker.C:
			fmt.Println("111")
		}
	}
}

var announcerStr = `
[{"ID":0,"TradeID":16639486,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640382764,"OpenPrice":1.1318,"ClosePrice":1.13072,"Comment":"13602387317_1","Profit":-1.08,"Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"TradeID":16639406,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640380965,"OpenPrice":1.13251,"ClosePrice":1.13072,"Comment":"13602387317_0","Profit":-1.79,"Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":16639372,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.02,"OpenTime":1640380842,"OpenPrice":1.13163,"ClosePrice":1.13188,"Comment":"13602387317_2","Profit":-0.5,"Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.02,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":16639297,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640379360,"OpenPrice":1.13104,"ClosePrice":1.13188,"Comment":"13602387317_1","Profit":-0.84,"Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":16639085,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640377076,"OpenPrice":1.13027,"ClosePrice":1.13188,"Comment":"13602387317_0","Profit":-1.61,"Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}}]`

var followerStr = `
[{"ID":0,"TradeID":15889024,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640408304,"OpenPrice":1.13122,"ClosePrice":1.13206,"Comment":"Copy BY #16638333,1000652,61057","Profit":0.84,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16638333}},{"ID":0,"TradeID":15888043,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640401938,"OpenPrice":1.1321,"ClosePrice":1.13206,"Comment":"Copy BY #16637870,1000652,61057","Profit":-0.04,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16637870}},{"ID":0,"Cmd":1,"TradeID":15887775,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640401190,"OpenPrice":1.13217,"ClosePrice":1.13247,"Comment":"Copy BY #16637795,1000652,61057","Profit":-0.3,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16637795}},{"ID":0,"TradeID":15887262,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640399839,"OpenPrice":1.13268,"ClosePrice":1.13206,"Comment":"Copy BY #16637544,1000652,61057","Profit":-0.62,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16637544}}]`

func TestV2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	relation := config.Relation{
		Announcer: "v1",
		Follower:  "v2",
	}

	var announcer []models.Order
	var follower []models.Order

	err := json.Unmarshal([]byte(announcerStr), &announcer)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal([]byte(followerStr), &follower)
	if err != nil {
		log.Fatalln(err)
	}

	var order []models.TwoWayOrder

	// 检查是否有 漏单
	for _, v := range announcer {
		exist := false
		for _, vv := range follower {
			if vv.Ex.TraderTradeID == v.TradeID {
				exist = true
			}

			fmt.Println("Tr: ", vv.Ex.TraderTradeID, "  td: ", v.TradeID)
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
		newOder = append(newOder, v)
	}

	if len(newOder) == 0 {
		return
	}

	fmt.Println("aaa-------")
	utils.Print(newOder)
}
