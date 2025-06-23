package pipeline

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/cmd"
)

const (
	KubeconfigPath = "/.kube/config"

	benchmarkNamespace    = "benchmark"
	bustCacheEnvVar       = "BUST_CACHE"
	kubeconfigEnvVar      = "KUBECONFIG"
	manifestFilename      = "/tmp/manifest.yaml"
	manifestFileExtension = "%s.yaml"
	podWaitDuration       = "15s"
	projectsDir           = "/projects"
	versionPlaceholder    = "$VERSION"
)

// Pipeline contains dagger objects that store the state of the pipeline and
// allow it to interact with the dagger container.
type Pipeline struct {
	container  *dagger.Container
	kubeconfig *dagger.File
	source     *dagger.Directory
}

func New(container *dagger.Container, source *dagger.Directory, kubeconfig *dagger.File) (*Pipeline, error) {
	return &Pipeline{
		container:  container,
		kubeconfig: kubeconfig,
		source:     source,
	}, nil
}

// Benchmark measures the sustainability footprint of CNCF projects.
// See README and docs directory for more details.
func (p *Pipeline) Benchmark(ctx context.Context,
	cncfProject,
	config,
	version,
	benchmarkJobURL string,
	benchmarkJobDurationMins int,
	prometheus_url string) (*dagger.Container, error) {
	if _, err := p.benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins); err != nil {
		log.Printf("benchmark failed: %v", err)
	}

	q, err := NewQuery(prometheus_url)

	if err != nil {
		log.Printf("failed to create prometheus query: %v", err)
		return nil, err
	}

	if results, err := p.computeBenchmarkingResults(ctx, q, benchmarkJobDurationMins); err != nil {
		log.Printf("failed to fetch metrics: %v", err)
		return nil, err
	} else {
		p.echo(ctx, "Benchmarking results:")
		results.WriteJSON(os.Stdout)
	}

	if _, err := p.delete(ctx, cncfProject, config, benchmarkJobURL); err != nil {
		return nil, err
	}

	return p.container, nil
}

// Terminal returns dagger interactive terminal configured with kubeconfig
// for trouble shooting.
func (p *Pipeline) Terminal(ctx context.Context) (*dagger.Container, error) {
	return p.withKubeconfig().Terminal(), nil
}

// benchmark creates the CNCF project resource, applies the benchmark tests and
// waits for the benchmark to complete.
func (p *Pipeline) benchmark(ctx context.Context,
	cncfProject,
	config,
	version,
	benchmarkJobURL string,
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	// Create CNCF project resources.
	if _, err := p.deploy(ctx, cncfProject, config, version); err != nil {
		return nil, err
	}

	// Create benchmark job resources.
	if _, err := p.exec(ctx, cmd.Apply(benchmarkJobURL)); err != nil {
		return nil, err
	}

	// Wait for pods to be ready.
	if _, err := p.exec(ctx, cmd.WaitForReadyPods(benchmarkNamespace)); err != nil {
		return nil, err
	}

	if _, err := p.echo(ctx, fmt.Sprintf("waiting %d minutes for benchmark to complete", benchmarkJobDurationMins)); err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(benchmarkJobDurationMins) * time.Minute)

	if _, err := p.echo(ctx, "benchmark complete"); err != nil {
		return nil, err
	}

	return p.container, nil
}

// delete cleans up the cluster at the end of the pipeline run.
func (p *Pipeline) delete(ctx context.Context, cncfProject, config, benchmarkJobURL string) (*dagger.Container, error) {
	// Delete benchmark job resources.
	if _, err := p.exec(ctx, cmd.Delete(benchmarkJobURL)); err != nil {
		log.Printf("failed to delete benchmark job: %v", err)
	}

	fileName, fileContents, err := p.getManifestFile(ctx, cncfProject, config, "")
	if err != nil {
		return nil, err
	}

	// Delete CNCF project resources.
	if _, err := p.execWithNewFile(ctx, fileName, fileContents, cmd.Delete(fileName)); err != nil {
		return nil, err
	}

	return p.container, nil
}

// deploy deploys the CNCF project by applying the flux manifest.
func (p *Pipeline) deploy(ctx context.Context, cncfProject, config, version string) (*dagger.Container, error) {
	fileName, fileContents, err := p.getManifestFile(ctx, cncfProject, config, version)
	if err != nil {
		return nil, err
	}

	if _, err = p.execWithNewFile(ctx, fileName, fileContents, cmd.Apply(fileName)); err != nil {
		return nil, err
	}

	// Allow time for pods to be created. This is needed because there is a
	// delay while flux reconciles the manifest. Without it the following
	// kubectl wait command will fail.
	if _, err := p.echo(ctx, fmt.Sprintf("waiting %s to let deployment complete", podWaitDuration)); err != nil {
		return nil, err
	}

	waitDuration, err := time.ParseDuration(podWaitDuration)
	if err != nil {
		return nil, err
	}
	time.Sleep(waitDuration)

	if _, err := p.exec(ctx, cmd.WaitForReadyPods(benchmarkNamespace)); err != nil {
		return nil, err
	}

	return p.container, nil
}

// echo outputs the message to stdout by running an echo command in the
// container. This is the dagger approach for logging to console output.
func (p *Pipeline) echo(ctx context.Context, msg string) (string, error) {
	// Bust cache to ensure commands are run.
	return p.container.WithEnvVariable(bustCacheEnvVar, time.Now().String()).
		WithExec(cmd.Echo(msg)).
		Stderr(ctx)
}

// exec configures the container with the kubeconfig and executes a command.
func (p *Pipeline) exec(ctx context.Context, args []string) (string, error) {
	return p.withKubeconfig().
		WithExec(args).
		Stdout(ctx)
}

// execWithDir is the same as exec but mounts a directory so manifests can be
// applied.
func (p *Pipeline) execWithDir(ctx context.Context, manifestPath string, args []string) (string, error) {
	dirPath := path.Dir(manifestPath)
	dir := p.source.Directory(dirPath)
	return p.withKubeconfig().
		WithDirectory(dirPath, dir).
		WithExec(args).
		Stdout(ctx)
}

// execWithNewFile is the same as exec but mounts a manifest file so it can be
// applied.
func (p *Pipeline) execWithNewFile(ctx context.Context, name, contents string, args []string) (string, error) {
	return p.withKubeconfig().
		WithNewFile(name, contents).
		WithExec(args).
		Stdout(ctx)
}

// getManifestFile gets the manifest and ensures the version is correct.
func (p *Pipeline) getManifestFile(ctx context.Context, project, config, version string) (string, string, error) {
	manifestPath := getManifestPath(project, config)
	manifest, err := p.source.File(manifestPath).Contents(ctx)
	if err != nil {
		return "", "", err
	}

	return manifestFilename, strings.ReplaceAll(manifest, versionPlaceholder, version), nil
}

// withKubeconfig adds the kubeconfig to the dagger container.
func (p *Pipeline) withKubeconfig() *dagger.Container {
	return p.container.
		// Bust cache with dynamic value to ensure layers are rebuilt and
		// commands are not cached by container engine.
		WithEnvVariable(bustCacheEnvVar, time.Now().String()).
		WithEnvVariable(kubeconfigEnvVar, KubeconfigPath).
		WithFile(KubeconfigPath, p.kubeconfig)
}

// getManifestPath get the correct path in the projects dir.
func getManifestPath(project, config string) string {
	if config == "" {
		config = project
	}
	return filepath.Join(projectsDir, project, fmt.Sprintf(manifestFileExtension, config))
}
