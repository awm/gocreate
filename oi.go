package gocreate

// Command is an interface representing a single OI command.
//
// Assemble generates the OI serial string that this command represents.
//
// Channel gets the Go channel to be used for returning responses to this command,
// if applicable.  This may return nil if no response channel is needed.
//
// Timeout is the time in ms to wait for a response if one is expected.  Positive
// values indicate a single response timeout, while negative values are used for an
// expected repeating (streaming) response.  The timeout is ignored if Channel()
// returns nil.
type Command interface {
    Assemble() []byte
    Channel() chan []byte
    Timeout() int
}

type simpleCommand struct {
    Opcode  byte
    Payload []byte
}

type baudCommand struct {
    Opcode  byte
    Payload []byte
    Rate    uint
}

// DemoNumber is a value indicating one of the Create's built in demo actions.
type DemoNumber byte

const (
    // DemoCover represents the "Cover" demo.
    DemoCover DemoNumber = iota
    // DemoCoverAndDock represents the "Cover and Dock" demo.
    DemoCoverAndDock
    // DemoSpotCover represents the "Spot Cover" demo.
    DemoSpotCover
    // DemoMouse represents the "Mouse" demo.
    DemoMouse
    // DemoFigureEight represents the "Drive Figure Eight" demo.
    DemoFigureEight
    // DemoWimp represents the "Wimp" demo.
    DemoWimp
    // DemoHome represents the "Home" demo.
    DemoHome
    // DemoTag represents the "Tag" demo.
    DemoTag
    // DemoPachelbel represents the "Pachelbel" demo.
    DemoPachelbel
    // DemoBanjo represents the "Banjo" demo.
    DemoBanjo
)

func (s *simpleCommand) Assemble() []byte {
    return append([]byte{s.Opcode}, s.Payload...)
}

func (s *simpleCommand) Channel() chan []byte {
    return nil
}

func (s *simpleCommand) Timeout() int {
    return 0
}

func (b *baudCommand) Assemble() []byte {
    return append([]byte{b.Opcode}, b.Payload...)
}

func (b *baudCommand) Channel() chan []byte {
    return nil
}

func (b *baudCommand) Timeout() int {
    return 0
}

// Start generates the "Start" command to initialize the OI.
func Start() Command {
    return &simpleCommand{Opcode: 128}
}

// Safe generates the "Safe" command to put the Create in Safe mode.
func Safe() Command {
    return &simpleCommand{Opcode: 131}
}

// Full generates the "Full" command to put the Create in Full control mode.
func Full() Command {
    return &simpleCommand{Opcode: 132}
}

// Baud generates the "Baud" command to change the connection baud rate.  The rate must
// be one of the supported OI baud rates:
//  300
//  600
//  1200
//  2400
//  4800
//  9600
//  14400
//  19200
//  28800
//  38400
//  57600
//  115200
func Baud(rate uint) Command {
    baudRates := map[uint][]byte{
        300:    {0},
        600:    {1},
        1200:   {2},
        2400:   {3},
        4800:   {4},
        9600:   {5},
        14400:  {6},
        19200:  {7},
        28800:  {8},
        38400:  {9},
        57600:  {10},
        115200: {11},
    }
    payload, ok := baudRates[rate]
    if !ok {
        return nil
    }

    return &baudCommand{Opcode: 129, Payload: payload, Rate: rate}
}

// Demo generates the "Demo" command to execute one of the built-in demos (see DemoNumber
// for the possible values).
func Demo(demo DemoNumber) Command {
    if demo < DemoCover || demo > DemoBanjo {
        return nil
    }

    return &simpleCommand{Opcode: 136, Payload: []byte{byte(demo)}}
}

// AbortDemo generates the "Demo" command to abort a currently executing built-in demo.
func AbortDemo() Command {
    return &simpleCommand{Opcode: 136, Payload: []byte{255}}
}

// Drive generates the "Drive" command tell the Create to move with a given velocity and
// turn radius.  The velocity must be in the range -500 to 500 mm/s, and the radius in
// the range -2000 to 2000 mm.
func Drive(velocity int16, radius int16) Command {
    if velocity < -500 || velocity > 500 {
        return nil
    }
    if radius < -2000 || radius > 2000 {
        return nil
    }

    payload := []byte{byte((velocity >> 8) & 0xFF), byte(velocity & 0xFF), byte((radius >> 8) & 0xFF), byte(radius & 0xFF)}
    return &simpleCommand{Opcode: 137, Payload: payload}
}

// DriveStraight generates the "Drive" command tell the Create to move with a given
// velocity in a straight line.  The velocity must be in the range -500 to 500 mm/s.
func DriveStraight(velocity int16) Command {
    if velocity < -500 || velocity > 500 {
        return nil
    }

    payload := []byte{byte((velocity >> 8) & 0xFF), byte(velocity & 0xFF), 0x80, 0x00}
    return &simpleCommand{Opcode: 137, Payload: payload}
}

// Spin generates the "Drive" command tell the Create to spin on the spot with a given
// velocity.  The velocity must be in the range -500 to 500 mm/s.
func Spin(velocity int16, clockwise bool) Command {
    if velocity < -500 || velocity > 500 {
        return nil
    }
    var high byte = 0x00
    var low byte = 0x01
    if clockwise {
        high, low = 0xFF, 0xFF
    }

    payload := []byte{byte((velocity >> 8) & 0xFF), byte(velocity & 0xFF), high, low}
    return &simpleCommand{Opcode: 137, Payload: payload}
}

// DriveDirect generates the "Drive Direct" command give the Create differential drive
// commands.  The left and right velocities must each be in the range of -500 to 500
// mm/s.
func DriveDirect(right int16, left int16) Command {
    if right < -500 || right > 500 {
        return nil
    }
    if left < -500 || left > 500 {
        return nil
    }

    payload := []byte{byte((right >> 8) & 0xFF), byte(right & 0xFF), byte((left >> 8) & 0xFF), byte(left & 0xFF)}
    return &simpleCommand{Opcode: 145, Payload: payload}
}

// Leds generates the "LEDs" command set the states of the onboard LEDs.
func Leds(advance bool, play bool, powerColour byte, powerIntensity byte) Command {
    var bits byte = 0x00
    if advance {
        bits |= (1 << 3)
    }
    if play {
        bits |= (1 << 1)
    }
    payload := []byte{bits, powerColour, powerIntensity}
    return &simpleCommand{Opcode: 139, Payload: payload}
}
