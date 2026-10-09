package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/NXWeb-Group/vnc-containers/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/google/uuid"
)

func main() {

	port := "2000"
	networkName := "vnc-network"

	// Create a new Docker client using the default configuration
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	err = utils.CreateImages(cli, "./docker")
	if err != nil {
		log.Fatalf("Failed to create Docker images: %v", err)
	}

	app := fiber.New()

	app.Use(static.New("./frontend/dist"))

	unusedContainer := map[string]bool{}
	var mutex sync.Mutex

	app.Get("/api/getImages", func(c fiber.Ctx) error {
		images, err := cli.ImageList(context.Background(), image.ListOptions{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to list images: " + err.Error())
		}

		names := []string{}
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if !strings.HasPrefix(tag, "vnc-") {
					continue
				}
				name := strings.TrimPrefix(tag, "vnc-")
				name = strings.TrimSuffix(name, ":latest")
				names = append(names, name)
			}
		}
		return c.JSON(names)
	})

	app.Get("/api/getContainers", func(c fiber.Ctx) error {
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to list containers: " + err.Error())
		}

		names := []string{}
		for _, container := range containers {
			for _, name := range container.Names {
				if strings.HasPrefix(strings.TrimPrefix(name, "/"), "vnc-instance-") {
					names = append(names, strings.TrimPrefix(name, "/"))
				}
			}
		}
		return c.JSON(names)
	})

	app.Post("/api/createContainer", func(c fiber.Ctx) error {
		var body struct{ Name string }
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body: " + err.Error())
		}

		id := uuid.NewString()
		containerName := "vnc-instance-" + id

		resp, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: "vnc-" + body.Name + ":latest",
		}, nil, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkName: {},
			},
		}, nil, containerName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create container: " + err.Error())
		}

		err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to start container: " + err.Error())
		}

		mutex.Lock()
		unusedContainer[containerName] = true
		mutex.Unlock()

		go utils.StartContainerTimer(cli, containerName, unusedContainer, &mutex, true)

		return c.JSON(fiber.Map{"id": id})
	})

	app.Post("/api/startContainer", func(c fiber.Ctx) error {
		var body struct{ Name string }
		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body: " + err.Error())
		}

		containerName := body.Name
		if !strings.HasPrefix(containerName, "vnc-instance-") {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid container name")
		}
		id := strings.TrimPrefix(containerName, "vnc-instance-")

		err = cli.ContainerStart(context.Background(), containerName, container.StartOptions{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to start container: " + err.Error())
		}

		mutex.Lock()
		unusedContainer[containerName] = true
		mutex.Unlock()

		go utils.StartContainerTimer(cli, containerName, unusedContainer, &mutex, false)

		return c.JSON(fiber.Map{"id": id})
	})

	app.Get("/websockify/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		log.Println("WebSocket connection established")
		utils.HandleWebSocket(c, id, cli, unusedContainer, &mutex)
	}))

	// handle ctrl c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server starting on", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server stopped with error: %v", err)
		}
	}()

	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v. Running cleanup...\n", sig)

	if err := app.Shutdown(); err != nil {
		fmt.Printf("Error shutting down Fiber: %v\n", err)
	}

	// app.Shutdown causes active ws containers to shutdown but this gets active conatiner names first
	// fix later
	utils.ShutdownContainers(cli)
	fmt.Println("Exiting")
}
