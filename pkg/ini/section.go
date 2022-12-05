package ini

import "log"

// TODO: Add support for merging/ignoring sections

// Section object provides contents of a particular section
type Section struct {
	values map[string]*Value
}

func newSection(sectName string) *Section {
	s := new(Section)
	s.values = make(map[string]*Value)
	return s
}

func (s *Section) addValue(valName, value string) *Value {
	val := newValue(valName, value)
	s.values[valName] = val
	return val
}

// Value returns named Value object from Section
func (s *Section) Value(valName string) *Value {
	val := s.ValueOptional(valName)
	if val == nil {
		log.Printf("INI.Section.Value: value '%s' not found", valName)
	}
	return val
}

// ValueOptional returns named Value object from Section
func (s *Section) ValueOptional(valName string) *Value {
	if s == nil {
		log.Printf("INI.Section.Value: looking for '%s' in undefined section", valName)
		return nil
	}
	if val, ok := s.values[valName]; ok {
		return val
	}
	return nil
}

// Values returns slice of value names
func (s *Section) Values() []string {
	keys := make([]string, len(s.values))

	i := 0
	for key := range s.values {
		keys[i] = key
		i++
	}
	return keys
}
