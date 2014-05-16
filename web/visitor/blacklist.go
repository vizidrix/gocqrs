package visitor_blacklist

import (
	"fmt"
	"github.com/vizidrix/gocqrs/cqrs"
	"github.com/vizidrix/gocqrs/web/visitor"
)

const DOMAIN_NAME string = "github.com/vizidrix/gocqrs/cqrs/web/visitor/blacklist"
const DOMAIN uint32 = 0x134F2B9D

const C_V1 uint32 = 0x80010000
const E_V1 uint32 = 0x00010000

const ( /* Commands */
	C_Blacklist        = cqrs.C(1, 1)
	C_RescindBlacklist = cqrs.C(1, 1)
)

const ( /* Events */
	E_Blacklisted        = cqrs.E(1, 1)
	E_BlacklistRescinded = cqrs.E(1, 1)
)

type Blacklist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewBlacklist(visitorId uint64, requestingUserId uint64) Blacklist {
	return Blacklist{
		CommandMemento: cqrs.NewCommand(visitor.DOMAIN, C_Blacklist, visitorId, -1),
		RequesterId:    requesterId,
	}
}

type RescindBlacklist struct {
	cqrs.CommandMemento
	RequesterId uint64 `json:"requestintUserId"`
}

func NewRescindBlacklist(visitorId uint64, requestingUserId uint64) RescindBlacklist {
	return RescindBlacklist{
		CommandMemento: cqrs.NewCommand(DOMAIN, C_RescindBlacklist, visitorId, -1),
		RequesterId:    requesterId,
	}
}
