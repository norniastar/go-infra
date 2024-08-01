package response

import errorcode "github.com/norniastar/go-infra/model/enum/error-code"

type API struct {
	Data  any             `json:"data"`
	Error errorcode.Value `json:"err"`
}
