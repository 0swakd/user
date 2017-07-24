package main

import (
    "os"
    "os/signal"
    "fmt"
    "log"
    "net/http"
    "syscall"
    "sync"
)

func init() {
    fmt.Fprintf(os.Stderr, "Checking server configuration\n")
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage : %s FILE\n")
        return
    }

    /* Meeeeh a context would be better but... */
    f, err := os.OpenFile(os.Args[1], os.O_RDONLY | os.O_CREATE, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

func main() {
    fmt.Fprintf(os.Stderr, "Initialising server\n")
    fmt.Fprintf(os.Stderr, "Loading routes\n")
    r := NewRouter()


    fmt.Fprintf(os.Stderr, "Configuring server\n")
    server := http.Server {
        Addr: ":8080",
        Handler: r,
    }

    fmt.Fprintf(os.Stderr, "Loading datas\n")
    StorageLoadUsers(os.Args[1])

    fmt.Fprintf(os.Stderr, "Configuring server graceful shutdown\n")
    interrupts := make(chan os.Signal, 1)

    signal.Notify(interrupts, syscall.SIGINT, syscall.SIGTERM)

    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        interrupt := <-interrupts
        defer wg.Done()

        fmt.Fprintf(os.Stderr, "Interrupt signal caught (%d), shutting down...\n", interrupt)
        err := server.Shutdown(nil)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error while shutting down %s\n", err)
        }

        fmt.Fprintf(os.Stderr, "Saving users\n")
        StorageSaveUsers(os.Args[1])
    }()

    fmt.Fprintf(os.Stderr, "Starting server\n")
    /* log.Fatal is screwing the data saving before gracefull shutdown */
    /* Should investigate this just to understand why */
    /*log.Fatal(*/server.ListenAndServe()/*)*/

    fmt.Fprintf(os.Stderr, "Main process : Waiting for data saving...\n")
    wg.Wait()

    fmt.Fprintf(os.Stderr, "Terminating server\n")
}

