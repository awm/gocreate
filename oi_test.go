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
