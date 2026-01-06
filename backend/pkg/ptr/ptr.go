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
