package hello_crypto_test

import (
	"fmt"
	"hello-go/hello_crypto"
)

func ExampleMd5Hex() {
	data := hello_crypto.Md5Hex("These pretzels are making me thirsty.")
	fmt.Printf("%s", data)
	// Output: b0804ec967f48520697662a204f5fe72
}
