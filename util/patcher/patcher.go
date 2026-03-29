package patcher

func Pick(updateMap map[string]any, allowed ...string) map[string]any {
	fields := make(map[string]bool, len(allowed))
	for _, f := range allowed {
		fields[f] = true
	}

	data := make(map[string]any)
	for key, value := range updateMap {
		if fields[key] {
			data[key] = value
		}
	}

	return data
}
