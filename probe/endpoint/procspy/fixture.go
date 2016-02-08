package procspy

type fixedConnIter []Connection

func (f *fixedConnIter) Next() *Connection {
	if len(*f) == 0 {
		return nil
	}

	car := (*f)[0]
	*f = (*f)[1:]

	return &car
}

// FixedScanner implements ConnectionScanner and uses constant Connection and
// ConnectionProcs.  It's designed to be used in tests.
type FixedScanner []Connection

// Connections implements ConnectionsScanner.Connections
func (s FixedScanner) Connections(_ bool) (ConnIter, error) {
	iter := fixedConnIter(s)
	return &iter, nil
}

// Stop implements ConnectionsScanner.Stop (dummy since there is no background work)
func (s FixedScanner) Stop() {}
