package main

import (
	"log"
	"net/http"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/miniplan"
	"github.com/gregoryv/miniplan/webui"
)

func main() {
	var (
		cli      = cmdline.NewBasicParser()
		bind     = cli.Option("-b, --bind").String(":9180")
		planfile = cli.Option("-f, --plan-file").String("plan.db")
	)
	cli.Parse()
	log.SetFlags(0)

	sys := miniplan.NewSystem()
	db, err := miniplan.NewPlanDB(planfile)
	if err != nil {
		log.Fatal(err)
	}
	sys.PlanDB = db
	ui := webui.NewUI(sys)
	if err := http.ListenAndServe(bind, ui); err != nil {
		log.Fatal(err)
	}
}
