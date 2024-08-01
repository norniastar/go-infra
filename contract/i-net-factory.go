package contract

import nettype "github.com/norniastar/go-infra/model/enum/net-type"

type INetFactory interface {
	Build(nettype.Value) INetService
}
