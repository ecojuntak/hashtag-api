package main

import (
	"log"
	"os"

	"github.com/ecojuntak/hashtag-api/cmds"
	"github.com/ecojuntak/hashtag-api/data"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "hashtag-api"
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:        "migrate",
			Description: "Run database migration",
			Action: func(c *cli.Context) error {
				err := data.RunMigration()
				if err != nil {
					log.Println(err)
				}
				return err
			},
		},
		{
			Name:        "start-amqp",
			Description: "Start listening to RabbitMQ Server",
			Action: func(c *cli.Context) error {
				cmds.StartRabbitMQ()
				return nil
			},
		},
		{
			Name:        "start-rest",
			Description: "Start listening to REST API",
			Action: func(c *cli.Context) error {
				cmds.StartREST()
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}
