package flags

import "fmt"

type stringArray []string

func (s *stringArray) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *stringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
}
