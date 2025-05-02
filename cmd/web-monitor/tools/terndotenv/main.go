package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

  fmt.Println(os.Getenv("DATABASE_NAME"))
  fmt.Println(os.Getenv("DATABASE_PASSWORD"))
  fmt.Println(os.Getenv("DATABASE_USER"))
  fmt.Println(os.Getenv("DATABASE_PORT"))
  fmt.Println(os.Getenv("DATABASE_HOST"))

	cmd := exec.Command("tern", "migrate", "--migrations", "./internal/store/pgstore/migrations", "--config", "./internal/store/pgstore/migrations/tern.conf")

	if err := cmd.Run(); err != nil {
    fmt.Println("Error running command:", err.Error())
		panic(err)
	}

}
