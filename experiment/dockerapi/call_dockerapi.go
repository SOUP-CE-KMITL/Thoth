package main

import (
    "fmt"

    "github.com/fsouza/go-dockerclient"
)

func main() {
    endpoint := "unix:///var/run/docker.sock"
    client, _ := docker.NewClient(endpoint)
    imgs, _ := client.ListVolume(docker.ListImagesOptions{})
    for _, img := range imgs {
        fmt.Println(img)
    }
}