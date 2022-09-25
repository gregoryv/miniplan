package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/miniplan"
	"github.com/gregoryv/miniplan/webui"
)

func main() {
	var (
		cli      = cmdline.NewBasicParser()
		bind     = cli.Option("-b, --bind").String("localhost:9180")
		planfile = cli.Option("-f, --plan-file").String("index.json")
	)
	cli.Parse()
	log.SetFlags(0)

	// open log file
	out, err := os.Create("mini.log")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	log.SetOutput(out)

	// create plan
	log.Print("create plan")
	dir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		dir = "."
	}

	plan := miniplan.NewPlan(dir, planfile)
	plan.Load()

	// init web user interface
	_ = webui.NewUI(plan)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}
