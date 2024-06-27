package features

import "github.com/xander42280/mpc/mpc_trade_test/common"

type Gateway interface {
	Post(path string, request interface{}) (*common.Response, error)
}
