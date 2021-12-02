package types

const (
	// UnspecifiedFsp is the unspecified fractional seconds part.
	UnspecifiedFsp = int8(-1)
	// DefaultFsp is the default digit of fractional seconds part.
	// MySQL use 0 as the default Fsp.
	DefaultFsp = int8(0)
)
