package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/urfave/cli"
	"github.com/williammartin/bozoci/image"
)

func main() {
	bozoci := cli.NewApp()
	bozoci.Name = "bozoci"
	bozoci.Commands = []cli.Command{
		CreateImageCommand,
	}

	if err := bozoci.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var CreateImageCommand = cli.Command{
	Name: "create-image",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "images-dir, d",
		},
		cli.StringFlag{
			Name: "image-uri, i",
		},
	},
	Action: func(ctx *cli.Context) error {
		provider := &image.Provider{
			ImagesDir: ctx.String("images-dir"),
		}

		imageURI, err := url.Parse(ctx.String("image-uri"))
		if err != nil {
			return err
		}

		rootfs, err := provider.Provide(ctx.Args().First(), imageURI)
		if err != nil {
			return err
		}

		fmt.Println(rootfs)
		return nil
	},
}
