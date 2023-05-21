package main

import (
        "dai-writer/routes"

        "github.com/gin-gonic/gin"
)

func main() {
        router := gin.Default()
        routes.AddPublics(router)
        routes.AddPrivates(router)
        router.Run(":5555")
}
