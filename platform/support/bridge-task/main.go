package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/neopilot-ai/neo/cmd/neo/mosaic/aws/appsync"
	"github.com/neopilot-ai/neo/cmd/neo/mosaic/aws/bridge"
	"github.com/neopilot-ai/neo/pkg/id"
)

var NEO_APP = os.Getenv("NEO_APP")
var NEO_STAGE = os.Getenv("NEO_STAGE")
var NEO_TASK_ID = os.Getenv("NEO_TASK_ID")
var NEO_REGION = os.Getenv("NEO_REGION")
var NEO_APPSYNC_HTTP = os.Getenv("NEO_APPSYNC_HTTP")
var NEO_APPSYNC_REALTIME = os.Getenv("NEO_APPSYNC_REALTIME")

var ENV_BLACKLIST = map[string]bool{
	"PATH": true,
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	workerID := id.Ascending()

	prefix := fmt.Sprintf("/neo/%s/%s", NEO_APP, NEO_STAGE)
	slog.Info("prefix", "value", prefix)
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(NEO_REGION))
	if err != nil {
		return err
	}

	conn, err := appsync.Dial(ctx, config, NEO_APPSYNC_HTTP, NEO_APPSYNC_REALTIME)
	if err != nil {
		return err
	}
	client := bridge.NewClient(ctx, conn, workerID, prefix+"/"+workerID)

	init := bridge.TaskStartBody{
		TaskID:      NEO_TASK_ID,
		Environment: []string{},
	}

	for _, e := range os.Environ() {
		key := strings.Split(e, "=")[0]
		if _, ok := ENV_BLACKLIST[key]; ok {
			continue
		}
		init.Environment = append(init.Environment, e)
	}
	creds, err := config.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	init.Environment = append(init.Environment, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID))
	init.Environment = append(init.Environment, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey))
	init.Environment = append(init.Environment, fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken))
	writer := client.NewWriter(bridge.MessageTaskStart, prefix+"/in")
	json.NewEncoder(writer).Encode(init)
	writer.Close()
	slog.Info("sent init")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-client.Read():
			if msg.Source != "dev" {
				continue
			}
			switch msg.Type {
			case bridge.MessagePing:
				slog.Info("got ping")
				continue
			case bridge.MessageTaskComplete:
				slog.Info("task complete")
				return nil
			}
		case <-time.After(time.Second * 10):
			slog.Info("timeout")
			return nil
		}
	}

}
