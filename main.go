package main

import (
	"fmt"
	"os"

	"encoding/csv"
	"github.com/urfave/cli"
	"log"
	"strconv"
	"strings"
	"time"
)

const STRING = "string"
const INT = "integer"
const FLOAT = "float"
const DATE = "date"

func main() {
	app := cli.NewApp()
	app.Name = "aschema"
	app.Usage = "Automatically detect schema and output parameters."
	app.Version = "0.1"

	app.Commands = []cli.Command{
		{
			Name:    "detect",
			Aliases: []string{"d"},
			Usage:   "Detect schema and output parameters.",
			Action: func(ctx *cli.Context) error {

				// csvの読み込み
				file, err := os.Open(ctx.Args().First())

				if err != nil {
					log.Fatal("Error:", err)
				}

				defer file.Close()

				reader := csv.NewReader(file)
				reader.LazyQuotes = true

				// ヘッダーの取得
				header, err := reader.Read()
				if err != nil {
					log.Fatal("Error:", err)
				}

				// 2行目移行の値を取得
				value, err := reader.Read()
				if err != nil {
					log.Fatal("Error:", err)
				}

				schema := []string{}
				for i, v := range header {
					s := fmt.Sprintf("%v:%v", v, getType(value[i]))
					schema = append(schema, s)
				}

				fmt.Println(strings.Join(schema, ","))

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func getType(str string) string {

	_, err := strconv.ParseFloat(str, 64)
	if err == nil && strings.Contains(str, ".") {
		return FLOAT
	}

	_, err = strconv.ParseInt(str, 10, 64)
	if err == nil {
		return INT
	}

	format := "2006-01-02 15:04:05"
	_, err = time.Parse(format, str)

	if err == nil {
		return DATE
	}

	return STRING
}
