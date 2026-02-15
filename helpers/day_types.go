package helpers

type dayType int

const (
	norm dayType = iota
	wknd
	off
	vac
	sic
)

func (d dayType) String() string {
	return [...]string{"norm", "wknd", "off", "vac", "sic"}[d]
}

func (d dayType) Next() dayType {
	return (d + 1) % 5
}

func ParseDayType(s string) dayType {
	switch s {
	case "wknd":
		return wknd
	case "off":
		return off
	case "vac":
		return vac
	case "sic":
		return sic
	default:
		return norm
	}
}
