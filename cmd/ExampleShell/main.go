package ExampleShell

import (
	"log"
)

func Execute() {
	log.Println("======= Start Example Shell ======= ")
	defer log.Println("======= Stop Example Shell ======= ")
}
