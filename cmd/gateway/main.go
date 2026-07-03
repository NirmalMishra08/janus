package gateway

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/config"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("not able to use config file")
		return
	}

	r:= router.SetupRouter()

   
	addr := fmt.Sprintf(":%s",cfg.PORT)
	http.ListenAndServe(addr, )


}
