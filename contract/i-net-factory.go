package contract

import nettype "github.com/infra-core/model/enum/net-type"

type INetFactory interface {
	Build(nettype.Value) INetService
}
