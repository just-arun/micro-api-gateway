package util

type array struct{}

func Array() *array {
	return &array{}
}

func (st array) Includes(items []string, cb func(item string, index int) bool) bool {
	for i, v := range items {
		if cb(v, i) {
			return true
		}
	}
	return false
}
