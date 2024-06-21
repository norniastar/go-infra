package response

import errorcode "github.com/infra-core/model/enum/error-code"

type API struct {
	Data  any             `json:"data"`
	Error errorcode.Value `json:"err"`
}
