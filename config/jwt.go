package config

import "os"

func JWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}
