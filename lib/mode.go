package lib

type mode int

const (
	c8 mode = iota
	sclr
	schr
)

func (m mode) String() string {
	return [...]string{"CHIP-8", "S-CHIP(LO-RES)", "S-CHIP(HI-RES)"}[m]
}
