package gocreate

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestStart(t *testing.T) {
    c := Start()
    assert.Equal(t, c.Assemble(), []byte{128}, "Assembled command string incorrect")
}

func TestSafe(t *testing.T) {
    c := Safe()
    assert.Equal(t, c.Assemble(), []byte{131}, "Assembled command string incorrect")
}

func TestFull(t *testing.T) {
    c := Full()
    assert.Equal(t, c.Assemble(), []byte{132}, "Assembled command string incorrect")
}

func TestBaud(t *testing.T) {
    valid := []uint{300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 28800, 38400, 57600, 115200}
    for i, b := range valid {
        c := Baud(b)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), []byte{129, byte(i)}, "Assembled command string for baud rate %u is incorrect", b)
        }
    }

    c := Baud(470)
    assert.Nil(t, c, "Expected creation of baud command to fail")
}

func TestDemo(t *testing.T) {
    for i := byte(DemoCover); i <= byte(DemoBanjo); i++ {
        c := Demo(DemoNumber(i))
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), []byte{136, i}, "Assembled command string for demo %u is incorrect", i)
        }
    }

    c := Demo(DemoBanjo + 100)
    assert.Nil(t, c, "Expected creation of too-large demo command to fail")
}

func TestAbortDemo(t *testing.T) {
    c := AbortDemo()
    assert.Equal(t, c.Assemble(), []byte{136, 255}, "Assembled command string incorrect")
}

type DrivePair struct {
    a   int16
    b   int16
    exp []byte
}

func TestDrive(t *testing.T) {
    valid := []DrivePair{
        {500, -2000, []byte{0x01, 0xF4, 0xF8, 0x30}},
        {-500, 2000, []byte{0xFE, 0x0C, 0x07, 0xD0}},
        {0, 0, []byte{0x00, 0x00, 0x00, 0x00}},
        {-27, 309, []byte{0xFF, 0xE5, 0x01, 0x35}},
        {400, 1100, []byte{0x01, 0x90, 0x04, 0x4C}},
        {400, -1100, []byte{0x01, 0x90, 0xFB, 0xB4}},
    }
    for _, p := range valid {
        c := Drive(p.a, p.b)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), append([]byte{137}, p.exp...), "Assembled command string for drive pair (%d, %d) is incorrect", p.a, p.b)
        }
    }

    c := Drive(501, 0)
    assert.Nil(t, c, "Expected creation of drive command with excessive positive velocity to fail")
    c = Drive(-501, 0)
    assert.Nil(t, c, "Expected creation of drive command with excessive negative velocity to fail")
    c = Drive(0, 2001)
    assert.Nil(t, c, "Expected creation of drive command with excessive positive radius to fail")
    c = Drive(0, -2001)
    assert.Nil(t, c, "Expected creation of drive command with excessive negative radius to fail")
}

func TestDriveStraight(t *testing.T) {
    valid := []DrivePair{
        {500, 0, []byte{0x01, 0xF4, 0x80, 0x00}},
        {-500, 0, []byte{0xFE, 0x0C, 0x80, 0x00}},
        {0, 0, []byte{0x00, 0x00, 0x80, 0x00}},
        {-27, 0, []byte{0xFF, 0xE5, 0x80, 0x00}},
        {400, 0, []byte{0x01, 0x90, 0x80, 0x00}},
    }
    for _, p := range valid {
        c := DriveStraight(p.a)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), append([]byte{137}, p.exp...), "Assembled command string for straight drive velocity %d is incorrect", p.a)
        }
    }

    c := DriveStraight(501)
    assert.Nil(t, c, "Expected creation of straight drive command with excessive positive velocity to fail")
    c = DriveStraight(-501)
    assert.Nil(t, c, "Expected creation of straight drive command with excessive negative velocity to fail")
}

func TestSpin(t *testing.T) {
    valid := []DrivePair{
        {200, 0, []byte{0x00, 0xC8, 0x00, 0x01}},
        {-200, 0, []byte{0xFF, 0x38, 0x00, 0x01}},
        {200, 1, []byte{0x00, 0xC8, 0xFF, 0xFF}},
        {-200, 1, []byte{0xFF, 0x38, 0xFF, 0xFF}},
    }
    for _, p := range valid {
        c := Spin(p.a, p.b != 0)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), append([]byte{137}, p.exp...), "Assembled command string for spin pair (%d, %t) is incorrect", p.a, p.b != 0)
        }
    }

    c := Spin(501, true)
    assert.Nil(t, c, "Expected creation of spin command with excessive positive velocity to fail")
    c = Spin(-501, true)
    assert.Nil(t, c, "Expected creation of spin command with excessive negative velocity to fail")
}

func TestDirectDrive(t *testing.T) {
    valid := []DrivePair{
        {200, -200, []byte{0x00, 0xC8, 0xFF, 0x38}},
        {0, 0, []byte{0x00, 0x00, 0x00, 0x00}},
        {200, 100, []byte{0x00, 0xC8, 0x00, 0x64}},
        {500, -500, []byte{0x01, 0xF4, 0xFE, 0x0C}},
    }
    for _, p := range valid {
        c := DriveDirect(p.a, p.b)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), append([]byte{145}, p.exp...), "Assembled command string for direct drive pair (%d, %d) is incorrect", p.a, p.b)
        }
    }

    c := DriveDirect(501, 0)
    assert.Nil(t, c, "Expected creation of direct drive command with excessive positive right velocity to fail")
    c = DriveDirect(-501, 0)
    assert.Nil(t, c, "Expected creation of direct drive command with excessive negative right velocity to fail")
    c = DriveDirect(0, 501)
    assert.Nil(t, c, "Expected creation of direct drive command with excessive positive left velocity to fail")
    c = DriveDirect(0, -501)
    assert.Nil(t, c, "Expected creation of direct drive command with excessive negative left velocity to fail")
}

func TestLeds(t *testing.T) {
    c := Leds(true, false, 128, 128)
    assert.Equal(t, c.Assemble(), []byte{139, 0x08, 128, 128}, "Assembled command string incorrect")
    c = Leds(false, true, 0, 255)
    assert.Equal(t, c.Assemble(), []byte{139, 0x02, 0, 255}, "Assembled command string incorrect")
}

func TestDigitalOutputs(t *testing.T) {
    for i := 0; i < 8; i++ {
        c := DigitalOutputs(byte(i))
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), []byte{147, byte(i)}, "Assembled command string for digital output value 0x%02X is incorrect", i)
        }
    }

    c := DigitalOutputs(0x08)
    assert.Nil(t, c, "Expected creation of digital output command with excessive value to fail")
}

type DriverTriplet struct {
    a   byte
    b   byte
    c   byte
}

func TestPwmLowSideDrivers(t *testing.T) {
    valid := []DriverTriplet{
        {0, 0, 0},
        {128, 128, 128},
        {1, 43, 86},
    }
    for _, v := range valid {
        c := PwmLowSideDrivers(v.a, v.b, v.c)
        if assert.NotNil(t, c) {
            assert.Equal(t, c.Assemble(), []byte{144, v.c, v.b, v.a}, "Assembled command string for PWM low side driver triplet (%u, %u, %u) is incorrect",
                v.a, v.b, v.c)
        }
    }

    c := PwmLowSideDrivers(129, 0, 0)
    assert.Nil(t, c, "Expected creation of PWM low side driver command with excessive driver 0 value to fail")
    c = PwmLowSideDrivers(0, 255, 0)
    assert.Nil(t, c, "Expected creation of PWM low side driver command with excessive driver 1 value to fail")
    c = PwmLowSideDrivers(0, 0, 129)
    assert.Nil(t, c, "Expected creation of PWM low side driver command with excessive driver 2 value to fail")
}

func TestSendIr(t *testing.T) {
    c := SendIr(157)
    assert.Equal(t, c.Assemble(), []byte{151, 157}, "Assembled command string incorrect")
}

func TestSong(t *testing.T) {
    valid := [][]Note{
        {{71, 33}, {74, 77}, {88, 100}, {121, 12}, {83, 90}, {72, 200}, {108, 150}},
        {{89, 25}, {72, 200}, {42, 44}, {36, 1}, {67, 90}, {55, 10}, {95, 200}, {90, 14}, {59, 110}, {106, 100}, {66, 80}, {33, 67}, {50, 222}, {94, 100},
            {127, 20}, {79, 50}},
        {{14, 64}, {66, 89}, {31, 200}, {121, 155}, {4, 200}},
    }
    for i, s := range valid {
        c := Song(byte(i), s)
        if assert.NotNil(t, c) {
            result := []byte{140, byte(i), byte(len(s))}
            for _, v := range s {
                result = append(result, v.Tone, v.Duration)
            }
            assert.Equal(t, c.Assemble(), result, "Assembled command string for Song %u is incorrect", i)
        }
    }

    c := Song(16, valid[0])
    assert.Nil(t, c, "Expected creation of Song command with excessive number to fail")
    c = Song(15, []Note{})
    assert.Nil(t, c, "Expected creation of Song command with empty list to fail")
    c = Song(15, append(valid[1], Note{100, 100}))
    assert.Nil(t, c, "Expected creation of Song command with more than 16 notes to fail")
}

func TestPlaySong(t *testing.T) {
    c := PlaySong(7)
    if assert.NotNil(t, c) {
        assert.Equal(t, c.Assemble(), []byte{141, 7}, "Assembled command string for song playback is incorrect")
    }

    c = PlaySong(20)
    assert.Nil(t, c, "Expected creation of PlaySong command with excessive number to fail")
}
