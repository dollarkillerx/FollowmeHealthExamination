package server

import (
	"fmt"
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

var announcer = `
[{"ID":0,"TradeID":16639486,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640382764,"OpenPrice":1.1318,"ClosePrice":1.13072,"Comment":"13602387317_1","Profit":-1.08,"
Commission":0,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"TradeID
":16639406,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640380965,"OpenPrice":1.13251,"ClosePrice":1.13072,"Comment":"13602387317_0","Profit":-1.79,"Commission":0,"E
x":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":166393
72,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.02,"OpenTime":1640380842,"OpenPrice":1.13163,"ClosePrice":1.13188,"Comment":"13602387317_2","Profit":-0.5,"Commission":0,"Ex":{"Stan
dardSymbol":"EUR/USD","StandardLots":0.02,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":16639297,"Symbo
l":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640379360,"OpenPrice":1.13104,"ClosePrice":1.13188,"Comment":"13602387317_1","Profit":-0.84,"Commission":0,"Ex":{"StandardSymb
ol":"EUR/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}},{"ID":0,"Cmd":1,"TradeID":16639085,"Symbol":"EURU
SD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640377076,"OpenPrice":1.13027,"ClosePrice":1.13188,"Comment":"13602387317_0","Profit":-1.61,"Commission":0,"Ex":{"StandardSymbol":"EUR
/USD","StandardLots":0.01,"FollowStatus":0,"TraderBrokerID":0,"TraderAccount":"","TraderName":"","TraderAccountIndex":0,"TraderTradeID":0}}]
`

var follower = `
[{"ID":0,"TradeID":15889024,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640408304,"OpenPrice":1.13122,"ClosePrice":1.13206,"Comment":"Copy BY #16638333,1000652,6105
7","Profit":0.84,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":
3,"TraderTradeID":16638333}},{"ID":0,"TradeID":15888043,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640401938,"OpenPrice":1.1321,"ClosePrice":1.13206,"Comment":"Cop
y BY #16637870,1000652,61057","Profit":-0.04,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"
Resse","TraderAccountIndex":3,"TraderTradeID":16637870}},{"ID":0,"Cmd":1,"TradeID":15887775,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1640401190,"OpenPrice":1.1321
7,"ClosePrice":1.13247,"Comment":"Copy BY #16637795,1000652,61057","Profit":-0.3,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":1,"TraderBrokerID":1000652,"T
raderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16637795}},{"ID":0,"TradeID":15887262,"Symbol":"EURUSD","ContractSize":100000,"Digits":5,"Volume":0.01,"OpenTime":1
640399839,"OpenPrice":1.13268,"ClosePrice":1.13206,"Comment":"Copy BY #16637544,1000652,61057","Profit":-0.62,"Commission":-0.04,"Ex":{"StandardSymbol":"EUR/USD","StandardLots":0.01,"FollowStatus":
1,"TraderBrokerID":1000652,"TraderAccount":"61057","TraderName":"Resse","TraderAccountIndex":3,"TraderTradeID":16637544}}]
`
