package types

// Set is for MySQL Set type.
type Set struct {
	Name string
}

// String implements fmt.Stringer interface.
func (e Set) String() string {
	return e.Name
}
