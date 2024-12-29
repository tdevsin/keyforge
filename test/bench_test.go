package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/tdevsin/keyforge/internal/proto"
)

func TestBenchmarkSetKey(t *testing.T) {
	_, close := runApp(t)
	defer close()
	hostname := "localhost:8080"
	t.Run("Benchmark Set Keys", func(t *testing.T) {
		report, err := runner.Run(
			"KeyService/SetKey",
			hostname,
			runner.WithProtoFile("../proto/keyforge.proto", []string{}),
			runner.WithInsecure(true),
			runner.WithConcurrency(1000),
			runner.WithRunDuration(time.Second*30),
			runner.WithDataProvider(func(*runner.CallData) ([]*dynamic.Message, error) {
				var dynamicMessages []*dynamic.Message
				for i := 0; i < 1000; i++ {
					protoMsg := &proto.SetKeyRequest{
						Key:   uuid.NewString(),
						Value: []byte("Hello world"),
					}
					dynamicMsg, err := dynamic.AsDynamicMessage(protoMsg)
					if err != nil {
						return nil, err
					}
					dynamicMessages = append(dynamicMessages, dynamicMsg)
				}

				return dynamicMessages, nil
			}),
		)
		if err != nil {
			t.Fatalf("Error running benchmark for %s: %v", "KeyService/SetKey", err)
		}

		printer := printer.ReportPrinter{
			Out:    os.Stdout,
			Report: report,
		}

		printer.Print("summary")
	})
}

func TestBenchmarkGetKey(t *testing.T) {
	_, close := runApp(t)
	defer close()
	hostname := "localhost:8080"
	t.Run("Benchmark Get Keys", func(t *testing.T) {

		conn := getGrpcConnection()
		var keys []string
		for i := 0; i < 1000; i++ {
			key := uuid.NewString()
			request := &proto.SetKeyRequest{
				Key:   key,
				Value: []byte("HelloWorld"),
			}
			client := proto.NewKeyServiceClient(conn)
			client.SetKey(context.Background(), request)
			keys = append(keys, uuid.NewString())
		}

		report, err := runner.Run(
			"KeyService/GetKey",
			hostname,
			runner.WithProtoFile("../proto/keyforge.proto", []string{}),
			runner.WithInsecure(true),
			runner.WithConcurrency(1000),
			runner.WithRunDuration(time.Second*30),
			runner.WithDataProvider(func(*runner.CallData) ([]*dynamic.Message, error) {
				var dynamicMessages []*dynamic.Message
				for i := 0; i < 1000; i++ {
					protoMsg := &proto.GetKeyRequest{
						Key: keys[i],
					}
					dynamicMsg, err := dynamic.AsDynamicMessage(protoMsg)
					if err != nil {
						return nil, err
					}
					dynamicMessages = append(dynamicMessages, dynamicMsg)
				}

				return dynamicMessages, nil
			}),
		)
		if err != nil {
			t.Fatalf("Error running benchmark for %s: %v", "KeyService/SetKey", err)
		}

		printer := printer.ReportPrinter{
			Out:    os.Stdout,
			Report: report,
		}

		printer.Print("summary")
	})
}
