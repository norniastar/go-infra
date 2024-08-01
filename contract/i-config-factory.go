package contract

import configgroup "github.com/norniastar/go-infra/model/enum/config-group"

type IConfigFactory interface {
	Build_(interface{}) IConfigService

	// Deprecated: Build_
	Build(group configgroup.Value) IConfigService
}
