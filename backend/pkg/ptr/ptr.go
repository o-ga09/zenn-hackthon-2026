package ptr

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func PtrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
