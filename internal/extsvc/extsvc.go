package extsvc

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewEmailService, NewProfileService)
