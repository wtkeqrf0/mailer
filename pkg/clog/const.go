package clog

type Level uint8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case 0:
		return `info`
	case 1:
		return `error`
	case 2:
		return `fatal`
	}
	return `<nil>`
}
