package pushswap

type Operation string

const (
	Invalid Operation = ""
	PA      Operation = "PA"
	PB      Operation = "PB"
	RA      Operation = "RA"
	RB      Operation = "RB"
	RR      Operation = "RR"
	RRA     Operation = "RRA"
	RRB     Operation = "RRB"
	RRR     Operation = "RRR"
	SA      Operation = "SA"
	SB      Operation = "SB"
	SS      Operation = "SS"
)
