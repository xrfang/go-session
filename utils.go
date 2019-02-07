package session

import "crypto/rand"

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func uuid(L int) string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	buf := make([]byte, L)
	rand.Read(buf)
	for i := 0; i < len(buf); i++ {
		p := buf[i] % 62
		buf[i] = charset[p]
	}
	return string(buf)
}
