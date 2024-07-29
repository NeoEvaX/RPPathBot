package bot

func splitMessage(s string) []string {
	var result []string
	var current string
	var inQuotes bool
	for _, char := range s {
		if char == '"' {
			inQuotes = !inQuotes
		} else if char == ' ' && !inQuotes {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}

func verifySetupGameMessage(message []string) bool {
	if message[0] != ".setupgame" {
		return false
	}
	if len(message) < 6 {
		return false
	}
	if message[2] != "GM" {
		return false
	}
	containsPlayers := func(message []string) bool {
		for _, a := range message {
			if a == "Players" {
				return true
			}
		}
		return false
	}
	if !containsPlayers(message) {
		return false
	}
	return true
}
