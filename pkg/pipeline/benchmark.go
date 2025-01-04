package pipeline

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/cmd"
)

const (
	benchmarkNamespace = "benchmark"
	falcoNamespace     = "falco" // TODO Remove when no longer used.
)

type Pipeline struct {
	container  *dagger.Container
	dir        *dagger.Directory
	kubeconfig *dagger.File
}

func New(container *dagger.Container, dir *dagger.Directory, kubeconfig *dagger.File) (*Pipeline, error) {
	return &Pipeline{
		container:  container,
		dir:        dir,
		kubeconfig: kubeconfig,
	}, nil
}

// Benchmark measures the sustainability footprint of CNCF projects.
func (p *Pipeline) Benchmark(ctx context.Context,
	cncfProject,
	config,
	version,
	benchmarkJobURL string,
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	_, err := p.benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins)
	if err != nil {
		log.Printf("benchmark failed: %v", err)
	}

	_, err = p.delete(ctx, cncfProject, config, benchmarkJobURL)
	if err != nil {
		return nil, err
	}

	return p.container, nil
}

func (p *Pipeline) Terminal(ctx context.Context) (*dagger.Container, error) {
	return p.withKubeconfig().Terminal(), nil
}

func (p *Pipeline) benchmark(ctx context.Context,
	cncfProject,
	config,
	version,
	benchmarkJobURL string,
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	// Create CNCF project resources.
	_, err := p.deploy(ctx, cncfProject, config, version)
	if err != nil {
		return nil, err
	}

	// Create benchmark job resources.
	_, err = p.exec(ctx, cmd.Apply(benchmarkJobURL))
	if err != nil {
		return nil, err
	}

	// Wait for pods to be ready.
	_, err = p.exec(ctx, cmd.WaitForNamespace(benchmarkNamespace))
	if err != nil {
		return nil, err
	}

	// TODO Remove once benchmark job resources created in benchmark namespace.
	_, err = p.exec(ctx, cmd.WaitForNamespace(falcoNamespace))
	if err != nil {
		return nil, err
	}

	p.echo(ctx, fmt.Sprintf("waiting %d minutes for benchmark to complete", benchmarkJobDurationMins))

	time.Sleep(time.Duration(benchmarkJobDurationMins) * time.Minute)

	p.echo(ctx, "benchmark complete")

	return p.container, nil
}

func (p *Pipeline) delete(ctx context.Context, cncfProject, config, benchmarkJobURL string) (*dagger.Container, error) {
	// Delete benchmark job resources.
	_, err := p.exec(ctx, cmd.Delete(benchmarkJobURL))
	if err != nil {
		log.Printf("failed to delete benchmark job: %v", err)
	}

	fileName, fileContents, err := p.getManifestFile(ctx, cncfProject, config, "")
	if err != nil {
		return nil, err
	}

	// Delete CNCF project resources.
	_, err = p.execWithNewFile(ctx, fileName, fileContents, cmd.Delete(fileName))
	if err != nil {
		return nil, err
	}

	return p.container, nil
}

func (p *Pipeline) deploy(ctx context.Context, cncfProject, config, version string) (*dagger.Container, error) {
	fileName, fileContents, err := p.getManifestFile(ctx, cncfProject, config, version)
	if err != nil {
		return nil, err
	}

	_, err = p.execWithNewFile(ctx, fileName, fileContents, cmd.Apply(fileName))
	if err != nil {
		return nil, err
	}

	// Allow time for pods to be created.
	time.Sleep(15 * time.Second)

	_, err = p.exec(ctx, cmd.WaitForNamespace(benchmarkNamespace))
	if err != nil {
		return nil, err
	}

	return p.container, nil
}

func (p *Pipeline) echo(ctx context.Context, msg string) (string, error) {
	// Bust cache to ensure commands are run.
	return p.container.WithEnvVariable("BUST_CACHE", time.Now().String()).
		WithExec([]string{"echo", fmt.Sprintf("'%s'", msg)}).
		Stdout(ctx)
}

func (p *Pipeline) exec(ctx context.Context, args []string) (string, error) {
	return p.withKubeconfig().
		WithExec(args).
		Stdout(ctx)
}

func (p *Pipeline) execWithDir(ctx context.Context, manifestPath string, args []string) (string, error) {
	dirPath := path.Dir(manifestPath)
	dir := p.dir.Directory(dirPath)
	return p.withKubeconfig().
		WithDirectory(dirPath, dir).
		WithExec(args).
		Stdout(ctx)
}

func (p *Pipeline) execWithNewFile(ctx context.Context, name, contents string, args []string) (string, error) {
	return p.withKubeconfig().
		WithNewFile(name, contents).
		WithExec(args).
		Stdout(ctx)
}

func (p *Pipeline) getManifestFile(ctx context.Context, project, config, version string) (string, string, error) {
	manifestPath := getManifestPath(project, config)
	manifest, err := p.dir.File(manifestPath).Contents(ctx)
	if err != nil {
		return "", "", err
	}

	return "/tmp/manifest.yaml", strings.ReplaceAll(manifest, "$VERSION", version), nil
}

func (p *Pipeline) withKubeconfig() *dagger.Container {
	return p.container.
		// Bust cache to ensure commands are run.
		WithEnvVariable("BUST_CACHE", time.Now().String()).
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithFile("/.kube/config", p.kubeconfig)
}

func getManifestPath(project, config string) string {
	if config == "" {
		config = project
	}
	return filepath.Join("/projects", project, fmt.Sprintf("%s.yaml", config))
}
