package pushswap

type Operation string

const (
	Invalid Operation = ""
	PA      Operation = "pa"
	PB      Operation = "pb"
	RA      Operation = "ra"
	RB      Operation = "rb"
	RR      Operation = "rr"
	RRA     Operation = "rra"
	RRB     Operation = "rrb"
	RRR     Operation = "rrr"
	SA      Operation = "sa"
	SB      Operation = "sb"
	SS      Operation = "ss"
)
