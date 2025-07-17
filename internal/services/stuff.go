package services

func StartsWithAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	b := s[0]
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
