package usecases

import "math/rand"

const charSet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM134567890"

func generateSlug(size int) string {
	var result string
	for range size {
		ind := rand.Intn(len(charSet))
		result += charSet[ind : ind+1]
	}
	return result
}
