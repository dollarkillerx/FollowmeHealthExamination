package models

// CurrentPosition 当前持仓
type CurrentPosition struct {
	Code int `json:"code"`
	Data struct {
		IsHiding       bool    `json:"IsHiding"`
		HasSubscribed  bool    `json:"HasSubscribed"`
		BrokerTimeZone int     `json:"BrokerTimeZone"`
		Total          int     `json:"Total"`
		TotalLots      float64 `json:"TotalLots"`
		TotalProfit    float64 `json:"TotalProfit"`
		Items          []Order `json:"Items"`
	} `json:"data"`
	Message string `json:"message"`
}

type Order struct {
	ID           int     `json:"ID"`
	Cmd          int     `json:"Cmd,omitempty"`
	TradeID      int     `json:"TradeID"`
	Symbol       string  `json:"Symbol"`
	ContractSize int     `json:"ContractSize"`
	Digits       int     `json:"Digits"`
	Volume       float64 `json:"Volume"`
	OpenTime     int     `json:"OpenTime"`
	OpenPrice    float64 `json:"OpenPrice"`
	ClosePrice   float64 `json:"ClosePrice"`
	Comment      string  `json:"Comment"`
	Profit       float64 `json:"Profit"`
	Commission   float64 `json:"Commission"`
	Ex           OrderEx `json:"Ex"`
}

type OrderEx struct {
	StandardSymbol     string  `json:"StandardSymbol"`
	StandardLots       float64 `json:"StandardLots"`
	FollowStatus       int     `json:"FollowStatus"`
	TraderBrokerID     int     `json:"TraderBrokerID"`
	TraderAccount      string  `json:"TraderAccount"`
	TraderName         string  `json:"TraderName"`
	TraderAccountIndex int     `json:"TraderAccountIndex"`
	TraderTradeID      int     `json:"TraderTradeID"` // 跟随订单号
}
