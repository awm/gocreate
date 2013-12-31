package gocreate

import (
    "fmt"
    "github.com/awm/goserial"
    "io"
    "time"
)

// Connection represents an active OI serial connection to a Create.
//
// Port is the name of the serial device that the OI connection is using.
//
// Baud is the  current baud rate of the serial connection.
type Connection struct {
    Port      string
    Baud      uint
    device    io.ReadWriteCloser
    sendQueue chan Command
}

// Connect open a new connection to a Create on the given serial port with the
// specified baud rate.  The Create must already be using the specified baud rate
// for the connection to function properly, though it can be changed later if necessary.
func Connect(port string, initialBaud uint) *Connection {
    config := &serial.Config{Name: port, Baud: int(initialBaud)}
    s, err := serial.OpenPort(config)
    if err != nil {
        fmt.Printf("Serial open error: %s\n", err.Error())
        return nil
    }

    conn := &Connection{Port: port, Baud: initialBaud, device: s, sendQueue: make(chan Command)}
    go conn.sender()
    return conn
}

func (c *Connection) sendData(data []byte) {
    if c.Baud == 115200 {
        for _, b := range data {
            time.Sleep(200 * time.Microsecond)
            c.device.Write([]byte{b})
        }
    } else {
        c.device.Write(data)
    }
}

func (c *Connection) sender() {
    for cmd := range c.sendQueue {
        data := cmd.Assemble()
        switch cmd := cmd.(type) {
        default:
            c.sendData(data)
        case *baudCommand:
            c.sendData(data)
            c.device.Close()

            c.Baud = cmd.Rate
            config := &serial.Config{Name: c.Port, Baud: int(c.Baud)}
            s, err := serial.OpenPort(config)
            if err != nil {
                panic(fmt.Sprintf("Failed to re-open serial port after baud rate change: %s", err.Error()))
            }
            c.device = s

            time.Sleep(100 * time.Millisecond)
        }
    }

    c.device.Close()
}

// Send transmits a single OI command to the connected Create.
func (c *Connection) Send(cmd Command) {
    c.sendQueue <- cmd
}

// SendMany transmits a sequence of OI commands to the connected Create.
func (c *Connection) SendMany(cmds []Command) {
    for _, cmd := range cmds {
        c.Send(cmd)
    }
}

// Close terminates the serial connection.
func (c *Connection) Close() {
    close(c.sendQueue)
}
