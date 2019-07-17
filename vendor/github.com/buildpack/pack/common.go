package pack

import (
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/buildpack/pack/builder"
	"github.com/buildpack/pack/style"
)

func (c *Client) parseTagReference(imageName string) (name.Reference, error) {
	if imageName == "" {
		return nil, errors.New("image is a required parameter")
	}
	if _, err := name.ParseReference(imageName, name.WeakValidation); err != nil {
		return nil, err
	}
	ref, err := name.NewTag(imageName, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a tag reference", imageName)
	}

	return ref, nil
}

func (c *Client) resolveRunImage(runImage, targetRegistry string, stackInfo builder.StackMetadata, additionalMirrors map[string][]string) string {
	if runImage != "" {
		c.logger.Debugf("Using provided run-image %s", style.Symbol(runImage))
		return runImage
	}
	runImageName := getBestRunMirror(
		targetRegistry,
		stackInfo.RunImage.Image,
		stackInfo.RunImage.Mirrors,
		additionalMirrors[stackInfo.RunImage.Image],
	)
	if runImageName == stackInfo.RunImage.Image {
		c.logger.Debugf("Selected run image %s", style.Symbol(runImageName))
	} else if contains(stackInfo.RunImage.Mirrors, runImageName) {
		c.logger.Debugf("Selected run image mirror %s", style.Symbol(runImageName))
	} else {
		c.logger.Debugf("Selected run image mirror %s from local config", style.Symbol(runImageName))
	}
	return runImageName
}

func contains(slc []string, v string) bool {
	for _, s := range slc {
		if s == v {
			return true
		}
	}
	return false
}

func getBestRunMirror(registry string, runImage string, mirrors []string, preferredMirrors []string) string {
	runImageList := append(preferredMirrors, append([]string{runImage}, mirrors...)...)
	for _, img := range runImageList {
		ref, err := name.ParseReference(img, name.WeakValidation)
		if err != nil {
			continue
		}
		if ref.Context().RegistryStr() == registry {
			return img
		}
	}

	if len(preferredMirrors) > 0 {
		return preferredMirrors[0]
	}

	return runImage
}
