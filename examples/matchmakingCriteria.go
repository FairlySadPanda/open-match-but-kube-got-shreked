package examples

import "fmt"

// search keys
const (
	GameModeKey          = "gameMode"
	RegionKey            = "region"
	TimeEnteringQueueKey = "timeEnteringQueue"
	EloKey               = "elo"
)

type GameMode int

const (
	EightBall GameMode = iota
	NineBall
	KoreanCarom
	JapaneseCarom
)

func (m GameMode) ToString() string {
	switch {
	case m == EightBall:
		return "eightBall"
	case m == NineBall:
		return "nineBall"
	case m == KoreanCarom:
		return "koreanCarom"
	case m == JapaneseCarom:
		return "japaneseCarom"
	default:
		// outright panic as something is very wrong
		panic(fmt.Sprintf("cannot convert gameMode to string: %v is not a mode", m))
	}
}

type Region int

const (
	// a really bad list of regions
	Europe Region = iota
	NorthAmerica
	SouthAmerica
	Asia
	Australasia
	Africa
)

func (r Region) ToString() string {
	switch {
	case r == Europe:
		return "europe"
	case r == NorthAmerica:
		return "northAmerica"
	case r == SouthAmerica:
		return "southAmerica"
	case r == Asia:
		return "asia"
	case r == Australasia:
		return "australasia"
	case r == Africa:
		return "africa"
	default:
		// outright panic as something is very wrong
		panic(fmt.Sprintf("cannot convert region to string: %v is not a region", r))
	}
}
