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
	"luuhai48/short/workers"

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

		cronjob := workers.CreateCleanupSessionCronjob()
		defer cronjob.Stop()
	}

	server.Get("", handlers.OptionalAuthMiddleware, handlers.HandleGetHomeIndex)

	server.Get("/signup", handlers.NoAuthGuard, handlers.HandleGetSignupIndex)
	server.Post("/signup", handlers.NoAuthGuard, handlers.HandlePostSignup)

	server.Get("/signin", handlers.NoAuthGuard, handlers.HandleGetSigninIndex)
	server.Post("/signin", handlers.NoAuthGuard, handlers.HandlePostSignin)

	server.Get("/signout", handlers.AuthMiddleware, handlers.HandleGetSignOut)

	server.Get("/short", handlers.AuthMiddleware, handlers.HandleShortIndex)
	server.Delete("short", handlers.AuthMiddleware, handlers.HandleDeleteShort)
	server.Get("/short/new", handlers.AuthMiddleware, handlers.HandleNewShortIndex)
	server.Post("/short/new", handlers.AuthMiddleware, handlers.HandlePostNewShort)

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

	server.Get("/:ID", handlers.HandleShortDetail)

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
