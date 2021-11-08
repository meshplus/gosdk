package types

// Enum is for MySQL enum type.
type Enum struct {
	Name string
}

// String implements fmt.Stringer interface.
func (e Enum) String() string {
	return e.Name
}
