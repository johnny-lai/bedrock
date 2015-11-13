package main

import (
	"github.com/codegangsta/cli"
	"github.com/johnny-lai/bedrock"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var version = "unset"

func bedrockRoot(c *cli.Context) (string, error) {
  root := c.GlobalString("base")
  log.Print(root)
  if root == "" {
    root = filepath.Dir(os.Args[0])
  }

  var err error
  root, err = filepath.Abs(root)
  if err != nil {
    return "", err
  }

  return root, nil
}

func executeTemplate(src string, dest string) error {
  tmpl, err := template.New(filepath.Base(src)).ParseFiles(src)
  if err != nil {
    return err
  }

  fd, err := os.Create(dest)
  if err != nil {
    return err
  }

  tc := bedrock.TemplateContext{}
  err = tmpl.Execute(fd, &tc)
  if err != nil {
    log.Fatal(err)
  }

  return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "bedrock"
	app.Version = version
	app.Usage = "A microservice structure for Go"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "base",
			Usage: "Location of bedrock",
      EnvVar: "BEDROCK_ROOT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "dump",
			Usage: "Reads the specified config file and prints the output",
			Action: func(c *cli.Context) {
				file := c.Args().First()

				tmpl, err := template.New(filepath.Base(file)).ParseFiles(file)
				if err != nil {
					log.Fatal(err)
					return
				}

				tc := bedrock.TemplateContext{}
				err = tmpl.Execute(os.Stdout, &tc)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
    {
			Name:  "generate",
			Usage: "<fixture>",
			Action: func(c *cli.Context) {
				fixture := c.Args().First()
        if fixture == "" {
          log.Fatal("Fixture must be specified")
        }

        target_path := fixture
        if len(c.Args()) > 1 {
          target_path = c.Args()[1]
        }
        target_path, err := filepath.Abs(target_path)
        if err != nil {
          log.Fatal(err)
        }

        root, _ := bedrockRoot(c)
        fixture_root := filepath.Join(root, "fixtures", fixture)
        finfo, err := os.Stat(fixture_root)
        if err != nil || !finfo.Mode().IsDir() {
          log.Fatalf("Unknown fixture %s specified. Expected %s to be a directory.", fixture, fixture_root)
        }


        log.Printf("Base Path: %s\n", root)

        generateFixture := (func (path string, info os.FileInfo, err error) error {
          basename := filepath.Base(path)
          if basename[0] == '.' {
            // Hidden file/directory. Ignore
            return filepath.SkipDir
          }

          rel_path := path[len(fixture_root):len(path)]
          out_path := filepath.Join(target_path, rel_path)
          if info.Mode().IsDir() {
            err := os.Mkdir(out_path, info.Mode())
            if err != nil && !os.IsExist(err)  {
              log.Fatal(err)
            }
          } else if info.Mode().IsRegular() {
            executeTemplate(path, out_path)
          }

          return nil
        })
        filepath.Walk(fixture_root, generateFixture)
			},
		},
	}

	app.Run(os.Args)
}
