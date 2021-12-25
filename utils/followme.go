package utils

import (
	"fmt"

	"github.com/dollarkillerx/FollowmeHealthExamination/pkg/models"
	"github.com/dollarkillerx/urllib"
)

// CurrentPosition 当前持仓
func CurrentPosition(uid string, token string, index string) (*models.CurrentPosition, error) {
	urli := fmt.Sprintf("https://www.followme.hk/api/v3/followtrade/agent/positions?uid=%s&index=%s", uid, index)

	var mo models.CurrentPosition
	err := urllib.Post(urli).RandUserAgent().
		SetJson([]byte(fmt.Sprintf(`{"OrderBy":0,"OrderField":4,"PageIndex":1,"PageSize":1000,"uid":%s,"index":13,"Type":0,"Page":1}`, uid))).
		SetHeader("expect-ct", `max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"`).
		SetHeader("server", "cloudflare").
		SetHeader("serverinfo", "www.followme.hk").
		SetHeader("cookie", token).
		FromJson(&mo)
	if err != nil {
		return nil, err
	}

	return &mo, nil
}
