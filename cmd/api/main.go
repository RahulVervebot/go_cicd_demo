package main
import "github.com/RahulVervebot/go_cicd_demo/internal/handlers"


func main() {
	mux := http.NewServeMux()
	handlers.Register(mux)

	addr := ":8080"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
