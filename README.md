# GoCreate

A Go wrapper for iRobot's [Create Open Interface (OI) Version 2](http://www.irobot.com/filelibrary/pdfs/hrd/create/Create%20Open%20Interface_v2.pdf).

## License

Licensed under the terms of the [Modified BSD License](http://opensource.org/licenses/BSD-3-Clause). See LICENSE.txt for details.

## Installation

In a GOPATH directory, run

    go get github.com/awm/gocreate

or, just clone this git repository.

## Usage

To use the library add an import statement to your Go file:

    import "github.com/awm/gocreate"

Then instatiate and use an OI connection in your code:

    conn := gocreate.Connect("/dev/tty.usbserial-A100eMiV", 57600)
    if conn == nil {
        return fmt.Errorf("Failed to open OI connection")
    }
    defer conn.Close()

    // Initialize and switch to "Full" mode
    cmds := []gocreate.Command{gocreate.Start(), gocreate.Full()}
    conn.SendMany(cmds)

    // Turn on all of the LEDs
    cmd := gocreate.Leds(true, true, 255, 255)
    conn.Send(cmd)

    // Do other things...

The file `connection_test.go` also demonstrates how to use the API.

## Documentation

Additional API documentation can be generated or viewed using the `godoc` tool:

    godoc github.com/awm/gocreate

## TODO

 * Improve tests (remove hardcoded serial port)
 * Implement remaining OI opcodes (currently only mode, demo, baud, and driving, and LED commands are implemented)
 * Improve error reporting/results
