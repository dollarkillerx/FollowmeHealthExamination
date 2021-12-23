package test

import (
	"github.com/dollarkillerx/FollowmeHealthExamination/pkg/models"
	"github.com/dollarkillerx/FollowmeHealthExamination/utils"
	"log"
	"testing"

	"github.com/dollarkillerx/urllib"
)

func TestFollowTest1(t *testing.T) {
	urli := "https://www.followme.hk/api/v3/followtrade/agent/positions?uid=188678&index=13"

	var mo models.CurrentPosition
	err := urllib.Post(urli).RandUserAgent().
		SetJson([]byte(`{"OrderBy":0,"OrderField":4,"PageIndex":1,"PageSize":1000,"uid":188678,"index":13,"Type":0,"Page":1}`)).
		SetHeader("expect-ct", `max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"`).
		SetHeader("server", "cloudflare").
		SetHeader("serverinfo", "www.followme.hk").
		SetHeader("cookie", "lang=zh-CN; theme=light; cookiesession1=678A3E0EABCDEFGHIJKLMNOPQRSV3329; HMF_CI=6209f3acb4e3557ed90250cc3b469235c015ca2f66085c75dc45601283bcc12594; HMY_JC=e5e8ca28520d67abf7f53f3c16a78832b300170ca208b1f8f9a60238a6d32563cc,; USER_TOKEN=7nUVbpLyYi4xqxPkSZu7VGeV50d2usBSjrcbs4FC9alReUw8oGKzMgkmgy83L58N9ut86cnIeWfKVLC7mGEbsA; AWSALBTG=gGuzJzkcCMVEya0mUUJpdrRPrDOUsgYVQRdj8LoMzLkHL/XnBPzwsx951VjKg3lWDz6xY2Tc4Hzo0Gsx8I27SM1gwVDaRfzLtvOdRuaEbJgQ+S07KBYavsMYxzYGdhLTgRUDNjmZwFmuEJ88AASevifUXZihPeEzN+FYMNIFncPeBSBLbB4=; AWSALBTGCORS=gGuzJzkcCMVEya0mUUJpdrRPrDOUsgYVQRdj8LoMzLkHL/XnBPzwsx951VjKg3lWDz6xY2Tc4Hzo0Gsx8I27SM1gwVDaRfzLtvOdRuaEbJgQ+S07KBYavsMYxzYGdhLTgRUDNjmZwFmuEJ88AASevifUXZihPeEzN+FYMNIFncPeBSBLbB4=").
		FromJson(&mo)
	if err != nil {
		log.Fatalln(err)
	}

	utils.Print(mo)
}

func TestFollowTest2(t *testing.T) {
	urli := "https://www.followme.hk/api/v3/followtrade/agent/positions?uid=784177&index=3"

	var mo models.CurrentPosition
	err := urllib.Post(urli).RandUserAgent().
		SetJson([]byte(`{"OrderBy":0,"OrderField":4,"PageIndex":1,"PageSize":1000,"uid":188678,"index":13,"Type":0,"Page":1}`)).
		SetHeader("expect-ct", `max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"`).
		SetHeader("server", "cloudflare").
		SetHeader("serverinfo", "www.followme.hk").
		SetHeader("cookie", "lang=zh-CN; theme=light; cookiesession1=678A3E0EABCDEFGHIJKLMNOPQRSV3329; HMF_CI=6209f3acb4e3557ed90250cc3b469235c015ca2f66085c75dc45601283bcc12594; HMY_JC=e5e8ca28520d67abf7f53f3c16a78832b300170ca208b1f8f9a60238a6d32563cc,; USER_TOKEN=7nUVbpLyYi4xqxPkSZu7VGeV50d2usBSjrcbs4FC9alReUw8oGKzMgkmgy83L58N9ut86cnIeWfKVLC7mGEbsA; AWSALBTG=gGuzJzkcCMVEya0mUUJpdrRPrDOUsgYVQRdj8LoMzLkHL/XnBPzwsx951VjKg3lWDz6xY2Tc4Hzo0Gsx8I27SM1gwVDaRfzLtvOdRuaEbJgQ+S07KBYavsMYxzYGdhLTgRUDNjmZwFmuEJ88AASevifUXZihPeEzN+FYMNIFncPeBSBLbB4=; AWSALBTGCORS=gGuzJzkcCMVEya0mUUJpdrRPrDOUsgYVQRdj8LoMzLkHL/XnBPzwsx951VjKg3lWDz6xY2Tc4Hzo0Gsx8I27SM1gwVDaRfzLtvOdRuaEbJgQ+S07KBYavsMYxzYGdhLTgRUDNjmZwFmuEJ88AASevifUXZihPeEzN+FYMNIFncPeBSBLbB4=").
		FromJson(&mo)
	if err != nil {
		log.Fatalln(err)
	}

	utils.Print(mo)
}
