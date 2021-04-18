package lsm303

import (
	"encoding/binary"
	"testing"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2ctest"
	"periph.io/x/periph/conn/mmr"
	"periph.io/x/periph/conn/physic"
)

var (
	magnetometerDatasheet = datasheetForMagnetometer(LSM303DLHC)
)

func TestNewMagnetometer(t *testing.T) {
	scenario := &i2ctest.Playback{
		Ops: []i2ctest.IO{
			// Read the chip ID (not a real ID, just a constant)
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.IRA_REG_M}, R: []byte{0b01001000}},
			// Read gain
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.CRB_REG_M}, R: []byte{0}},
			// Write new gain
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.CRB_REG_M, uint8(MAGNETOMETER_GAIN_4_0) << 5}, R: []byte{}},
			// Write new rate
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.CRA_REG_M, (uint8(MAGNETOMETER_RATE_30) << 2) | 0b10000000}, R: []byte{}},
		},
	}

	_, err := NewMagnetometer(scenario)

	if err != nil {
		t.Fatal(err)
	}
}

func TestMagnetometerSense(t *testing.T) {
	scenario := &i2ctest.Playback{
		Ops: []i2ctest.IO{
			// Read registers
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_X_L_M}, R: []byte{0}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_X_H_M}, R: []byte{1}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_Y_L_M}, R: []byte{100}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_Y_H_M}, R: []byte{0}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_Z_L_M}, R: []byte{0xff}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.OUT_Z_H_M}, R: []byte{0xff}},
		},
	}

	magnetometer := &Magnetometer{
		mmr: mmr.Dev8{
			Conn:  &i2c.Dev{Bus: scenario, Addr: magnetometerDatasheet.ADDRESS},
			Order: binary.BigEndian,
		},
		gain: MAGNETOMETER_GAIN_4_0,
		rate: MAGNETOMETER_RATE_30,
	}

	x, y, z, err := magnetometer.SenseRaw()
	if err != nil {
		t.Fatal(err)
	}

	if x != 256 {
		t.Fatal("Bad x")
	}
	if y != 100 {
		t.Fatal("Bad y")
	}
	if z != -1 {
		t.Fatal("Bad z")
	}
}

func TestGetTemperature(t *testing.T) {
	scenario := &i2ctest.Playback{
		Ops: []i2ctest.IO{
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_H_M}, R: []byte{0}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_L_M}, R: []byte{0}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_H_M}, R: []byte{0}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_L_M}, R: []byte{0b10000000}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_H_M}, R: []byte{0b11111111}},
			{Addr: magnetometerDatasheet.ADDRESS, W: []byte{magnetometerDatasheet.TEMP_OUT_L_M}, R: []byte{0b10000000}},
		},
	}

	magnetometer := &Magnetometer{
		mmr: mmr.Dev8{
			Conn:  &i2c.Dev{Bus: scenario, Addr: magnetometerDatasheet.ADDRESS},
			Order: binary.BigEndian,
		},
		gain: MAGNETOMETER_GAIN_4_0,
		rate: MAGNETOMETER_RATE_30,
	}

	temperature, err := magnetometer.SenseRelativeTemperature()
	if err != nil {
		t.Fatal(err)
	}
	const offset = 20
	if temperature != physic.ZeroCelsius+offset*physic.Celsius {
		t.Fatal("Not 0 C")
	}

	temperature, err = magnetometer.SenseRelativeTemperature()
	if err != nil {
		t.Fatal(err)
	}
	if temperature != physic.ZeroCelsius+physic.Celsius+offset*physic.Celsius {
		t.Fatal("Not 1 C")
	}

	temperature, err = magnetometer.SenseRelativeTemperature()
	if err != nil {
		t.Fatal(err)
	}
	if temperature != physic.ZeroCelsius-physic.Celsius+offset*physic.Celsius {
		t.Fatal("Not -1 C")
	}
}
