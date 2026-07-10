package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("your-chosen-admin-password"), bcrypt.DefaultCost)
	fmt.Println(string(hash))
}
