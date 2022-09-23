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
		bind     = cli.Option("-b, --bind").String(":9180")
		planfile = cli.Option("-f, --plan-file").String("plan.db")
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

	// open database
	log.Print("open ", planfile)
	db, err := miniplan.NewPlanDB(planfile)
	if err != nil {
		log.Fatal(err)
	}

	// create system
	log.Print("create system")
	sys := miniplan.NewSystem()
	sys.SetDatabase(db)

	// init web user interface
	_ = webui.NewUI(sys)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}
