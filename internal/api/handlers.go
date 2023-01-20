// Package api provides the web interface
package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/foroozf001/broker-service/internal/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// JSONPayload is the type for json posted to this api
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// LogViaGRPC takes a JSON payload and logs it via gRPC
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, r, err, http.StatusBadRequest)
		return
	}

	loggerServiceURI, exists := os.LookupEnv("LOGGER_SERVICE_URI")
	if !exists {
		log.Println("missing logger service uri")
		return
	}

	connection, err := grpc.Dial(loggerServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		_ = app.errorJSON(w, r, err, http.StatusInternalServerError)
		return
	}
	defer connection.Close()

	c := pb.NewLogServiceClient(connection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &pb.LogRequest{
		LogEntry: &pb.Log{
			Name: requestPayload.Name,
			Data: requestPayload.Data,
		},
	})
	if err != nil {
		_ = app.errorJSON(w, r, err, http.StatusInternalServerError)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged event " + requestPayload.Name

	_ = app.writeJSON(w, r, http.StatusAccepted, payload)
}
