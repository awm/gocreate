package gocreate

import (
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestBlink(t *testing.T) {
    conn := Connect("/dev/tty.usbserial-A100eMiV", 57600)
    if !assert.NotNil(t, conn) {
        return
    }
    defer conn.Close()

    var cmds []Command
    cmd := Start()
    assert.NotNil(t, cmd)
    cmds = append(cmds, cmd)

    cmd = Full()
    assert.NotNil(t, cmd)
    cmds = append(cmds, cmd)

    cmd = Leds(true, true, 255, 255)
    assert.NotNil(t, cmd)
    cmds = append(cmds, cmd)

    conn.SendMany(cmds)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, true, 255, 255)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 255, 255)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 128, 255)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 0, 255)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 0, 170)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 0, 85)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 0, 0)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
    time.Sleep(2 * time.Second)

    cmd = Leds(false, false, 255, 255)
    assert.NotNil(t, cmd)
    conn.Send(cmd)
}
