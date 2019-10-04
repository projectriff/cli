package build

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"

	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/cache"
	"github.com/buildpack/pack/logging"
	"github.com/buildpack/pack/style"
)

// PlatformAPIVersion is the current Platform API Version supported by this version of pack.
const PlatformAPIVersion = "0.1"

type Lifecycle struct {
	builder      *builder.Builder
	logger       logging.Logger
	docker       *client.Client
	appPath      string
	appOnce      *sync.Once
	httpProxy    string
	httpsProxy   string
	noProxy      string
	LayersVolume string
	AppVolume    string
}

type Cache interface {
	Name() string
	Clear(context.Context) error
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func NewLifecycle(docker *client.Client, logger logging.Logger) *Lifecycle {
	return &Lifecycle{logger: logger, docker: docker}
}

type LifecycleOptions struct {
	AppPath    string
	Image      name.Reference
	Builder    *builder.Builder
	RunImage   string
	ClearCache bool
	Publish    bool
	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
}

func (l *Lifecycle) Execute(ctx context.Context, opts LifecycleOptions) error {
	l.Setup(opts)
	defer l.Cleanup()

	buildCache := cache.NewVolumeCache(opts.Image, "build", l.docker)
	launchCache := cache.NewVolumeCache(opts.Image, "launch", l.docker)
	l.logger.Debugf("Using build cache volume %s", style.Symbol(buildCache.Name()))

	if opts.ClearCache {
		if err := buildCache.Clear(ctx); err != nil {
			return errors.Wrap(err, "clearing build cache")
		}
		l.logger.Debugf("Build cache %s cleared", style.Symbol(buildCache.Name()))
	}

	l.logger.Debug(style.Step("DETECTING"))
	if err := l.Detect(ctx); err != nil {
		return err
	}

	l.logger.Debug(style.Step("RESTORING"))
	if opts.ClearCache {
		l.logger.Debug("Skipping 'restore' due to clearing cache")
	} else {
		if err := l.Restore(ctx, buildCache.Name()); err != nil {
			return err
		}
	}

	l.logger.Debug(style.Step("ANALYZING"))
	if err := l.Analyze(ctx, opts.Image.Name(), opts.Publish, opts.ClearCache); err != nil {
		return err
	}

	l.logger.Debug(style.Step("BUILDING"))
	if err := l.Build(ctx); err != nil {
		return err
	}

	l.logger.Debug(style.Step("EXPORTING"))
	launchCacheName := launchCache.Name()
	if err := l.Export(ctx, opts.Image.Name(), opts.RunImage, opts.Publish, launchCacheName); err != nil {
		return err
	}

	l.logger.Debug(style.Step("CACHING"))
	if err := l.Cache(ctx, buildCache.Name()); err != nil {
		return err
	}
	return nil
}

func (l *Lifecycle) Setup(opts LifecycleOptions) {
	l.LayersVolume = "pack-layers-" + randString(10)
	l.AppVolume = "pack-app-" + randString(10)
	l.appPath = opts.AppPath
	l.appOnce = &sync.Once{}
	l.builder = opts.Builder
	l.httpProxy = opts.HTTPProxy
	l.httpsProxy = opts.HTTPSProxy
	l.noProxy = opts.NoProxy
}

func (l *Lifecycle) Cleanup() error {
	var reterr error
	if err := l.docker.VolumeRemove(context.Background(), l.LayersVolume, true); err != nil {
		reterr = errors.Wrapf(err, "failed to clean up layers volume %s", l.LayersVolume)
	}
	if err := l.docker.VolumeRemove(context.Background(), l.AppVolume, true); err != nil {
		reterr = errors.Wrapf(err, "failed to clean up app volume %s", l.AppVolume)
	}
	return reterr
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a' + byte(rand.Intn(26))
	}
	return string(b)
}
