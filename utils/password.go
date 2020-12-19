package utils

import "math/rand"

//GeneratePassowrd -
func GeneratePassowrd(length int, withspec bool) string {
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"

	rand.Seed(NowJST().Time().UnixNano())
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits
	if withspec {
		all += specials
	}

	off := 1
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	if withspec {
		buf[1] = specials[rand.Intn(len(specials))]
		off = 2
	}
	alen := len(all)
	for i := off; i < length; i++ {
		buf[i] = all[rand.Intn(alen)]
	}
	rand.Shuffle(length, func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	return string(buf)
}
