package image

import (
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/containers/image/image"
	"github.com/containers/image/transports"
	"github.com/containers/image/types"

	_ "github.com/containers/image/docker"
)

type Provider struct {
	ImagesDir string
}

func (p *Provider) Provide(imageID string, imageURL *url.URL) (string, error) {
	imageSource, sourcedImage, err := getSourceAndImage(imageURL, &types.SystemContext{
		OSChoice:                    "linux",
		DockerInsecureSkipTLSVerify: true,
	})
	if err != nil {
		return "", err
	}
	defer func() {
		imageSource.Close()
	}()

	// get all layer metadata
	layerInfos := sourcedImage.LayerInfos()
	imageDir := filepath.Join(p.ImagesDir, imageID)
	// unpack each layer onto disk
	for _, layerInfo := range layerInfos {
		blobStream, err := getBlobStream(imageSource, layerInfo)
		if err != nil {
			return "", err
		}

		if err := os.MkdirAll(imageDir, 0777); err != nil {
			return "", err
		}

		tarCmd := exec.Command("tar", "-z", "-p", "-x", "-C", imageDir)
		tarCmd.Stdin = blobStream
		if err := tarCmd.Run(); err != nil {
			return "", err
		}

		blobStream.Close()
	}

	return imageDir, nil
}

func getSourceAndImage(imageURL *url.URL, systemContext *types.SystemContext) (types.ImageSource, types.Image, error) {
	ref, err := reference(imageURL)
	if err != nil {
		return nil, nil, err
	}

	imageSource, err := ref.NewImageSource(systemContext)
	if err != nil {
		return nil, nil, err
	}

	sourcedImage, err := image.FromSource(systemContext, imageSource)
	if err != nil {
		imageSource.Close()
		return nil, nil, err
	}

	return imageSource, sourcedImage, nil
}

func getBlobStream(imageSource types.ImageSource, layer types.BlobInfo) (io.ReadCloser, error) {
	blobStream, _, err := imageSource.GetBlob(layer)
	if err != nil {
		return nil, err
	}

	return blobStream, err
}

func reference(imageURL *url.URL) (types.ImageReference, error) {
	transport := transports.Get(imageURL.Scheme)

	refString := "/"
	if imageURL.Host != "" {
		refString += "/" + imageURL.Host
	}
	refString += imageURL.Path

	ref, err := transport.ParseReference(refString)
	if err != nil {
		return nil, err
	}

	return ref, nil
}
