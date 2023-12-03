package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func MatchPasswordToHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
// main functions in shared directory are used as example usage.
// to run main function, rename the package as main and then run: go run ./path/to/this/util/*.go
func main() {
	password := "secret"
	passwordHash, err := Hash(password)
	if err != nil {
		fmt.Printf("error hashing password: %v", err)
		return
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", passwordHash)

	match := MatchPasswordToHash(password, passwordHash)
	fmt.Println("Match:   ", match)
}
