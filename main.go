package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

// Config Config.
type Config struct {
	Filelayout FileLayoutConfig
}

// FileLayoutConfig FileLayoutConfig.
type FileLayoutConfig struct {
	LayoutType []LayoutType `toml:"type"`
}

// LayoutType tyoe of fixed-length file layout.
type LayoutType struct {
	ID     string `toml:"id"`
	Length []int  `toml:"length"`
}

func main() {

	app := cli.NewApp()

	app.Name = "cvfv"
	app.Usage = "Convert Fixed-length file into variable-length file"
	app.Version = "0.0.1"

	var layout string
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "layout, l",
			Usage:       "Fixed-length file layout",
			Aliases:     []string{"l"},
			Destination: &layout,
		},
	}

	app.Action = func(context *cli.Context) error {

		// option check
		if layout == "" {
			fmt.Fprintln(os.Stderr, "ERROR : --layout option is Required.")
			os.Exit(1)
		}

		// config file check
		pathOfConfigFile := os.Getenv("HOME") + "/.config/cvfv/config.toml"
		lengthList, err := readLengthListByToml(pathOfConfigFile, layout)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR :", err)
			os.Exit(1)
		}

		switch context.Args().Len() {
		case 0:
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fmt.Println(convCSV(scanner.Text(), lengthList))
			}
		case 1:
			in, err := os.Open(context.Args().First())
			if err != nil {
				panic(err)
			}
			defer in.Close()
			scanner := bufio.NewScanner(in)
			for scanner.Scan() {
				fmt.Println(convCSV(scanner.Text(), lengthList))
			}
		default:
			fmt.Fprintln(os.Stderr, "ERROR : Invalid number of argument.")
			os.Exit(1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func readLengthListByToml(tomlPath string, ltype string) (retList []int, ferr error) {

	// file exist check
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		msg := "file doesn't exit. [" + tomlPath + "]\n"
		return nil, errors.New(msg)
	}

	// settings read
	var config Config
	_, confErr := toml.DecodeFile(tomlPath, &config)
	if confErr != nil {
		panic(confErr)
	}
	//	lengthList := getLengthListByLayoutType(config.Filelayout.LayoutType, layout)

	for _, item := range config.Filelayout.LayoutType {
		if item.ID == ltype {
			return item.Length, nil
		}
	}

	msg := "layout doesn't exit in setting's toml"
	return nil, errors.New(msg)

}

func convCSV(str string, layout []int) (ret string) {
	const camma = ","

	if len(str) == 0 {
		return ""
	}

	from := 0
	for _, digit := range layout {
		ret += string([]rune(str)[from:from+digit]) + camma
		from += digit
	}

	return strings.TrimRight(ret, camma)
}
