package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"
	"flag"

	cloudevents "github.com/cloudevents/sdk-go"

	"github.com/knative-sample/weather-store/pkg/controller"
	"github.com/knative-sample/weather-store/pkg/kncloudevents"
	"github.com/knative-sample/weather-store/pkg/utils/logs"
)

func dispatch(ctx context.Context, event cloudevents.Event) {
	fmt.Printf(event.String())
	controller.StoreWeather()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "ok")
}
func main() {
	flag.Parse()
	logs.InitLogs()
	defer logs.FlushLogs()
	go func() {
		http.HandleFunc("/health", handler)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8022"
		}
		http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}()

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), dispatch))

}
