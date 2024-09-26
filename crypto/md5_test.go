package crypto_test

import (
	"fmt"
	"hello-go/crypto"
)

func ExampleMd5Hex() {
	data := crypto.Md5Hex("These pretzels are making me thirsty.")
	fmt.Printf("%s", data)
	// Output: b0804ec967f48520697662a204f5fe72
}
