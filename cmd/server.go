package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"

	"luuhai48/short/db"
	"luuhai48/short/handlers"
	"luuhai48/short/static"
	"luuhai48/short/utils"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
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

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if !fiber.IsChild() {
		if err := migrateDB(); err != nil {
			log.Fatal(err)
		}
	}

	server.Get("", handlers.HandleGetHomeIndex)

	server.Get("/signup", handlers.HandleGetSignupIndex)
	server.Post("/signup", handlers.HandlePostSignup)

	server.Get("/signin", handlers.HandleGetSigninIndex)
	server.Post("/signin", handlers.HandlePostSignin)

	server.Get("/_url", handlers.AuthMiddleware, handlers.HandleUrlHomeIndex)

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
