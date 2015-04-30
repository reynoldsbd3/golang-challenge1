package drum

import (
	"fmt"
	"os"
)

func main() {

	pattern, err := DecodeFile(os.Args[1])
    if err != nil {
        panic(err)
    }

	fmt.Println(pattern)
}
