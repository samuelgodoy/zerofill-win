package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
    "time"

    "github.com/go-ole/go-ole"
    "github.com/go-ole/go-ole/oleutil"
    "github.com/urfave/cli/v2"
)

var throughput int64 = 10 * 1024 * 1024 // Throughput inicial em bytes por segundo.

func listDisks() {
    ole.CoInitialize(0)
    defer ole.CoUninitialize()

    unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
    defer unknown.Release()

    wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
    defer wmi.Release()

    serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
    service := serviceRaw.ToIDispatch()
    defer service.Release()

    resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Win32_DiskDrive")
    result := resultRaw.ToIDispatch()
    defer result.Release()

    countVar, _ := oleutil.GetProperty(result, "Count")
    count := int(countVar.Val)

    for i := 0; i < count; i++ {
        itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
        item := itemRaw.ToIDispatch()
        defer item.Release()

        model, _ := oleutil.GetProperty(item, "Model")
        deviceID, _ := oleutil.GetProperty(item, "DeviceID")
        fmt.Printf("Device %d: %s - %s\n", i, model.ToString(), deviceID.ToString())
    }
}

func zeroFill(devicePath string) {
    fmt.Printf("Starting zero fill on %s with initial throughput %d MB/s...\n", devicePath, throughput/(1024*1024))

    file, err := os.OpenFile(devicePath, os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open device %s for writing: %v", devicePath, err)
    }
    defer file.Close()

    var totalWritten int64 = 0

    go func() {
        reader := bufio.NewReader(os.Stdin)
        fmt.Println("Press '+' to increase throughput by 100KB/s, '-' to decrease.")
        for {
            char, _, err := reader.ReadRune()
            if err != nil {
                fmt.Println("Error reading from keyboard:", err)
                return
            }
            switch char {
            case '+':
                throughput += 100 * 1024 // Aumenta 100 KB/s.
                fmt.Printf("\nThroughput adjusted to: %d KB/s", throughput/1024)
            case '-':
                throughput -= 100 * 1024 // Diminui 100 KB/s.
                if throughput < 0 {
                    throughput = 0
                }
                fmt.Printf("\nThroughput adjusted to: %d KB/s", throughput/1024)
            }
        }
    }()

    for {
        // Quando o throughput é zero, pausa o zero fill.
        if throughput == 0 {
            fmt.Println("\nThroughput is zero. Zero fill paused. Press '+' to resume.")
            for throughput == 0 {
                time.Sleep(1 * time.Second) // Espera ativa até que o throughput seja aumentado.
            }
        } else {
            // Ajusta o tamanho do buffer com base no throughput atual.
            bufferSize := throughput
            if bufferSize > 1024*1024 { // Limita o tamanho do buffer a 1MB.
                bufferSize = 1024 * 1024
            }

            buffer := make([]byte, bufferSize)
            written, err := file.Write(buffer)
            if err != nil {
                if err == io.EOF {
                    fmt.Println("\nZero fill completed successfully.")
                    return
                }
                log.Fatalf("Failed to write to device: %v", err)
            }
            totalWritten += int64(written)
            fmt.Printf("\rTotal written: %d MB", totalWritten/(1024*1024))
            
            // Calcula o tempo de espera baseado no throughput desejado para evitar spikes.
            if bufferSize > 0 {
                time.Sleep(time.Second * time.Duration(bufferSize) / time.Duration(throughput))
            }
        }
    }
}

func main() {
    app := &cli.App{
        Name:  "ZeroFiller",
        Usage: "Zero fill your drives with adjustable throughput.",
        Commands: []*cli.Command{
            {
                Name:  "list",
                Usage: "List all physical drives",
                Action: func(c *cli.Context) error {
                    listDisks()
                    return nil
                },
            },
            {
                Name:  "zero",
                Usage: "Perform zero fill on a specified drive",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:     "device",
                        Aliases:  []string{"d"},
                        Required: true,
                        Usage:    "The device path to zero fill",
                    },
                },
                Action: func(c *cli.Context) error {
                    devicePath := c.String("device")
                    zeroFill(devicePath)
                    return nil
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
