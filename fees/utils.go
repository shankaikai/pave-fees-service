package fees

func contains(arr []string, str string) bool {
	for _, s := range arr {
			if s == str {
					return true
			}
	}
	return false
}