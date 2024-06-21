package response

import errorcode "github.com/norniastar/infra-core/model/enum/error-code"

type API struct {
	Data  any             `json:"data"`
	Error errorcode.Value `json:"err"`
}
