package RANode

import "time"

type Arg struct {
	RequestMsg int
	TimeStamp  time.Time
	Version    int
}

type Reply struct {
	ReplyMsg int
}
