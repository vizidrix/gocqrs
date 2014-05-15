package web

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
)

var DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web"
var DOMAIN uint32 = 0xD937B693

var ( /* Commands */
	C_StartServer			= cqrs.C(1, 1)
	C_StopServer 			= cqrs.C(1, 2)
	C_HandleHttpRequest    	= cqrs.C(1, 3)
)

var ( /* Events */
	E_ServerStarted			= cqrs.E(1, 1)
	E_ServerStopped			= cqrs.E(1, 2)
	E_HttpRequestHandled    = cqrs.E(1, 4)
)
