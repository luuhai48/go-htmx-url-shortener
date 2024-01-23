package main

import (
	"log"
	"log/slog"
	"luuhai48/short/static"
	"luuhai48/short/utils"
	"luuhai48/short/views"
	"os"
	"os/signal"

	"github.com/a-h/templ"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/urfave/cli/v2"
)

func startServer(ctx *cli.Context) error {
	server := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	server.Use(
		recover.New(recover.Config{EnableStackTrace: true}),
	)

	server.Get("", adaptor.HTTPHandler(templ.Handler(views.Index())))

	server.Use(
		"/static",
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		filesystem.New(filesystem.Config{
			Root:   static.StaticFilesFS(),
			MaxAge: 604800,
		}),
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		slog.Info("Gracefully shutting down...")
		server.Shutdown()
	}()

	if err := server.Listen(utils.MustGetEnv("ADDR", "0.0.0.0:3333")); err != nil {
		log.Fatal(err)
	}

	return nil
}
