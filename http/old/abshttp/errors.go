
type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

type ErrorType uint

const ()
