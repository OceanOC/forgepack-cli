package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v3"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cmd := &cli.Command{
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "file",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "embed",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "outFolder",
				Value: pwd,
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.StringArg("file") == "" {
				fmt.Printf("Not enough arguments. (atleast 1 needed)\nPress ENTER to exit\n")
				if runtime.GOOS == "windows" {
					fmt.Println("HINT: Try dragging and dropping the zip file onto the executable")
				}
				if !cmd.Bool("embed") {
					fmt.Scanln()
				}
				os.Exit(1)
			}

			var cffile CFManifest
			OpenZIP(cmd.StringArg("file"), cmd.String("outFolder"), &cffile)

			fmt.Println("All done! \nPress ENTER to exit")
			fmt.Scanln()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
