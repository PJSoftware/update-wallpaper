package ini

// TODO: Add support for blank/empty values

// TODO: Perhaps the calling code needs the ability to specifically state the
//  default value at the point of the call, rather than by a separate call?

// Value stores the actual named value, and provides conversion functions
type Value struct {
	strValue string
}

func newValue(valName, value string) *Value {
	v := new(Value)
	v.strValue = value
	return v
}
