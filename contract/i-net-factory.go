package contract

import nettype "github.com/norniastar/infra-core/model/enum/net-type"

type INetFactory interface {
	Build(nettype.Value) INetService
}
