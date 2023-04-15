package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// usage: go run ./build/ci/test/main.go
//
// This is a proof of concept. It's too slow to be used in codespaces.
func main() {
	fmt.Println("ci/test")
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Question: how to store cache on actions?
	goModCache := client.CacheVolume("gomod")
	goBinCache := client.CacheVolume("gobin")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	hostSourceDir := client.Host().Directory(wd, dagger.HostDirectoryOpts{
		Exclude: []string{},
	})

	exec := func(cmd ...string) []string {
		return cmd
	}

	const workspace = "/workspace"

	golang := client.Container().
		From("mcr.microsoft.com/devcontainers/go:0-1.20-bullseye").
		WithMountedCache("/go/pkg/mod", goModCache).
		WithMountedCache("/go/bin", goBinCache).
		WithMountedDirectory(workspace, hostSourceDir).
		WithWorkdir(workspace)

	_, err = golang.
		WithExec(exec("make")).
		WithExec(exec("git", "diff", "--exit-code")).
		Stderr(ctx)
	if err != nil {
		panic(err)
	}
}
