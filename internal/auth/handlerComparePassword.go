package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(hash []byte, password string) error{
        return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
