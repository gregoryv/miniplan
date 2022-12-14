package main

import (
	"fmt"
	"io/ioutil"
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
		bind     = cli.Option("-b, --bind").String(":9180")
		planfile = cli.Option("-f, --plan-file").String("index.json")
		logfile  = cli.Option("-l, --log-file").String("")
	)
	cli.Parse()
	log.SetFlags(0)

	if logfile != "" {
		out, err := os.Create(logfile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("create", logfile)
		defer out.Close()
		log.SetOutput(out)
	}

	// create plan
	if _, err := os.Stat(planfile); err != nil {
		if err := ioutil.WriteFile(planfile, []byte("{}"), 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("create", planfile)
	}

	plan := miniplan.NewPlan(planfile)
	plan.Load()

	// init web user interface
	fmt.Printf("listen on %s\n", bind)
	_ = webui.NewUI(plan)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}
