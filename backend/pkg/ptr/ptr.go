package ptr

import "database/sql"

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

func Int64ToPtr(i int64) *int {
	if i == 0 {
		return nil
	}
	i2 := int(i)
	return &i2
}

func BoolToPtr(b bool) *bool {
	return &b
}

func PtrToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func PtrStringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func PtrIntToNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func PtrBoolToNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}
