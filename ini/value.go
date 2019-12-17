package ini

import (
	"log"
	"strconv"
)

// Value stores the actual named value, and provides conversion functions
type Value struct {
	key      string
	strValue string
}

func newValue(valName, value string) *Value {
	v := new(Value)
	v.key = valName
	v.strValue = value
	return v
}

// AsString returns value as string
func (v *Value) AsString(def string, blankPermitted bool) string {
	if v == nil {
		log.Printf("INI.Section.Value.AsString: value undefined; returning default")
		return def
	}
	if v.strValue == "" {
		if !blankPermitted {
			return def
		}
	}
	return v.strValue
}

// AsBool returns value as bool
func (v *Value) AsBool(def bool) bool {
	if v == nil {
		log.Printf("INI.Section.Value.AsBool: value undefined; returning default")
		return def
	}
	if v.strValue == "" {
		return def
	}
	val, err := strconv.ParseBool(v.strValue)
	if err != nil {
		log.Printf("INI.Section.Value.AsBool: Unable to convert '%s' value '%s' to bool; returning '%v'", v.key, v.strValue, def)
		log.Print(err)
		return def
	}
	return val
}

// AsInt returns value as int (or default if unable to convert)
func (v *Value) AsInt(def int) int {
	if v == nil {
		log.Printf("INI.Section.Value.AsInt: value undefined; returning default")
		return def
	}
	if v.strValue == "" {
		return def
	}
	val, err := strconv.ParseInt(v.strValue, 10, 0)
	if err != nil {
		log.Printf("INI.Section.Value.AsInt: Unable to convert '%s' value '%s' to int64; returning '%d'", v.key, v.strValue, def)
		log.Print(err)
		return def
	}
	return int(val)
}

// AsInt64 returns value as int64 (or default if unable to convert)
func (v *Value) AsInt64(def int64) int64 {
	if v == nil {
		log.Printf("INI.Section.Value.AsInt64: value undefined; returning default")
		return def
	}
	if v.strValue == "" {
		return def
	}
	val, err := strconv.ParseInt(v.strValue, 10, 64)
	if err != nil {
		log.Printf("INI.Section.Value.AsInt64: Unable to convert '%s' value '%s' to int64; returning '%d'", v.key, v.strValue, def)
		log.Print(err)
		return def
	}
	return val
}
