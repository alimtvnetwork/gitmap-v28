package appfault

// Exception is an alias for any, representing a captured panic value.
type Exception = any

// Block represents a Try/Catch/Finally paradigm in Go.
type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

// Do executes the Try block. If a panic occurs, it is recovered and passed to Catch.
// Finally is executed unconditionally at the end of the operation, regardless of panics.
func (it Block) Do() {
	if it.Finally != nil {
		defer it.Finally()
	}

	if it.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				it.Catch(r)
			}
		}()
	}

	it.Try()
}
