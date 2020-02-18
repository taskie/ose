package ose

func rejectEmpty(ss []string) []string {
	results := make([]string, 0)
	for _, s := range ss {
		if s != "" {
			results = append(results, s)
		}
	}
	return results
}
