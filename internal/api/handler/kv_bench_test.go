package handler

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

var c *config.Config

func getConfig() (*config.Config, func()) {
	f := func() {
		os.RemoveAll("testdir")
	}
	if c != nil {
		return c, f
	}
	os.RemoveAll("./testdir")
	os.Mkdir("./testdir", 0755)

	l := logger.GetLogger()
	l.Logger = zap.NewNop()

	d := storage.GetDatabaseInstance(l, "./testdir")
	c = &config.Config{
		RootDir: "testdir",
		Logger:  l,
		Db:      d,
	}
	return c, f
}
func BenchmarkSetKey(b *testing.B) {
	config, _ := getConfig()

	h := KVHandler{Conf: config}
	b.SetParallelism(512)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := h.SetKey(context.Background(), &proto.SetKeyRequest{
				Key:   uuid.NewString(),
				Value: []byte(uuid.NewString()),
			})
			if err != nil {
				b.Logf("SetKey failed: %v", err)
			}
		}
	})
}

func BenchmarkGetKey(b *testing.B) {
	config, _ := getConfig()
	h := KVHandler{
		Conf: config,
	}
	keys := make([]string, 512)
	for i := 0; i < 512; i++ {
		keys[i] = uuid.NewString()
		h.SetKey(context.Background(), &proto.SetKeyRequest{
			Key:   keys[i],
			Value: []byte("HelloWorld"),
		})
	}
	b.SetParallelism(512)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randomNumber := rand.Intn(512)
			_, err := h.GetKey(context.Background(), &proto.GetKeyRequest{
				Key: keys[randomNumber],
			})
			if err != nil {
				b.Logf("Error: %v", err)
			}
		}
	})
}
