package utils

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func createTarArchive(srcDir string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	tw := tar.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer tw.Close()

		err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Create tar header
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}

			// Set the header name to the relative path
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}
			header.Name = relPath

			// Write header
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// Write file content if it's a regular file
			if info.Mode().IsRegular() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(tw, file)
				return err
			}

			return nil
		})

		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	return pr, nil
}

func buildImage(cli *client.Client, ctx context.Context, dockerfilePath string, imageName string) error {
	buildContext, err := createTarArchive(dockerfilePath)
	if err != nil {
		return fmt.Errorf("Failed to create build context: %v", err)
	}
	// Build options
	buildOptions := build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"vnc-" + imageName + ":latest"},
		Remove:     true,
	}

	// Build the image
	response, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return fmt.Errorf("Failed to build image: %v", err)
	}
	defer response.Body.Close()

	// Stream the build output
	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read build output: %v", err)
	}
	return nil
}

func CreateImages(cli *client.Client, dockerfilePath string) error {

	entries, err := os.ReadDir(dockerfilePath)
	if err != nil {
		return fmt.Errorf("Failed to read Dockerfile directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err := buildImage(cli, context.Background(), filepath.Join(dockerfilePath, entry.Name()), entry.Name())
			if err != nil {
				return fmt.Errorf("Failed to build Docker images: %v", err)
			}
		}
	}

	log.Println("Docker images built successfully!")
	return nil
}

func StartContainerTimer(cli *client.Client, containerName string, unusedContainer map[string]bool, mutex *sync.Mutex) {
	time.Sleep(1 * time.Minute)
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := unusedContainer[containerName]; exists {
		go cleanupContainer(cli, containerName)
		delete(unusedContainer, containerName)
	}
}

func cleanupContainer(cli *client.Client, containerName string) {
	if cli == nil {
		log.Printf("Docker client is nil, cannot cleanup container %s", containerName)
		return
	}

	// Stop the container
	err := stopContainer(cli, containerName)
	if err != nil {
		log.Printf("Failed to stop container %s: %s", containerName, err)
	} else {
		log.Printf("Container %s stopped successfully", containerName)
	}

	// Remove the container
	err = cli.ContainerRemove(context.Background(), containerName, container.RemoveOptions{
		Force: true, // Force removal even if running
	})
	if err != nil {
		log.Printf("Failed to remove container %s: %s", containerName, err)
	} else {
		log.Printf("Container %s removed successfully", containerName)
	}
}

func stopContainer(cli *client.Client, containerName string) error {
	timeout := 10 // 10 seconds timeout
	ctx := context.Background()
	err := cli.ContainerStop(ctx, containerName, container.StopOptions{
		Timeout: &timeout,
	})
	return err
}

func ShutdownContainers(cli *client.Client) {
	if cli == nil {
		log.Printf("Docker client is nil, cannot shutdown Chrome containers")
		return
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list Chrome containers: %s", err)
		return
	}

	for _, item := range containers {
		for _, name := range item.Names {
			if strings.HasPrefix(strings.TrimPrefix(name, "/"), "vnc-instance-") {
				stopContainer(cli, item.ID)
				break
			}
		}
	}
}
