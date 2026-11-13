// 程序入口。组装和生命周期在 internal/app。
package main

import (
	"log"

	"gin-enterprise/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
