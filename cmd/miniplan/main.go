package main

import (
	"log"
	"net/http"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/miniplan"
)

func main() {
	var (
		cli      = cmdline.NewBasicParser()
		bind     = cli.Option("-b, --bind").String(":9180")
		planfile = cli.Option("-f, --plan-file").String("plan.db")
	)
	log.SetFlags(0)

	sys := miniplan.NewSystem()
	db, err := miniplan.NewPlanDB(planfile)
	if err != nil {
		log.Fatal(err)
	}
	sys.PlanDB = db
	if err := http.ListenAndServe(bind, sys); err != nil {
		log.Fatal(err)
	}
}
