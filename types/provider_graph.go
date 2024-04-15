package types

type GraphRecord struct {
	// Values contains all the values in the record.
	Values []interface{}
	// Keys contains names of the values in the record.
	// Should not be modified. Same instance is used for all records within the same result.
	Keys []string
}
