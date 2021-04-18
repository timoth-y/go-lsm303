package lsm303

import (
	"encoding/binary"
	"fmt"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
	"periph.io/x/periph/conn/physic"
)

// This is a handle to the LSM303 magnetometer sensor.
type Magnetometer struct {
	mmr  mmr.Dev8
	sensorType SensorType
	datasheet *MagnetometerDatasheet
	addr *uint16
	rate MagnetometerRate
	gain MagnetometerGain
}

// New magnetometer opens a handle to an LSM303 magnetometer sensor.
func NewMagnetometer(bus i2c.Bus, opts ...MagnetometerOption) (*Magnetometer, error) {
	device := &Magnetometer{
		sensorType: LSM303DLHC,
		gain: MAGNETOMETER_GAIN_4_0,
		rate: MAGNETOMETER_RATE_30,
	}

	for i := range opts {
		opts[i].Apply(device)
	}

	if device.datasheet == nil {
		device.datasheet = datasheetForMagnetometer(device.sensorType)
	}

	if device.addr == nil {
		device.addr = &device.datasheet.ADDRESS
	}

	device.mmr = mmr.Dev8{
		Conn: &i2c.Dev{Bus: bus, Addr: *device.addr},
		// I don't think we ever access more than 1 byte at once, so
		// this is irrelevant
		Order: binary.BigEndian,
	}

	// Enable the magnetometer
	err := device.mmr.WriteUint8(device.datasheet.MR_REG_M, 0x00)
	if err != nil {
		return nil, err
	}

	// Validate sensor
	if chipId, err := device.mmr.ReadUint8(device.datasheet.WHO_AM_I_M); err != nil {
		return nil, err
	} else if chipId != device.datasheet.CHIP_ID {
		return nil, fmt.Errorf("no %s detected", device.sensorType)
	}

	// Init magnetometer configuration
	device.SetGain(device.gain)
	device.SetRate(device.rate)

	return device, nil
}



func (m *Magnetometer) SenseRaw() (int16, int16, int16, error) {
	xLow, err := m.mmr.ReadUint8(m.datasheet.OUT_X_L_M)
	if err != nil {
		return 0, 0, 0, err
	}
	xHigh, err := m.mmr.ReadUint8(m.datasheet.OUT_X_H_M)
	if err != nil {
		return 0, 0, 0, err
	}
	yLow, err := m.mmr.ReadUint8(m.datasheet.OUT_Y_L_M)
	if err != nil {
		return 0, 0, 0, err
	}
	yHigh, err := m.mmr.ReadUint8(m.datasheet.OUT_Y_H_M)
	if err != nil {
		return 0, 0, 0, err
	}
	zLow, err := m.mmr.ReadUint8(m.datasheet.OUT_Z_L_M)
	if err != nil {
		return 0, 0, 0, err
	}
	zHigh, err := m.mmr.ReadUint8(m.datasheet.OUT_Z_H_M)
	if err != nil {
		return 0, 0, 0, err
	}

	xValue := int16(((uint16(xHigh)) << 8) + uint16(xLow))
	yValue := int16(((uint16(yHigh)) << 8) + uint16(yLow))
	zValue := int16(((uint16(zHigh)) << 8) + uint16(zLow))

	return xValue, yValue, zValue, nil
}

func (m *Magnetometer) SetRate(mode MagnetometerRate) error {
	const bits = 3
	const shift = 2

	// The only bit in here that matters is the bit 7, temperature
	// enabled, so just always set it to 1
	previous := uint8(0)
	data := uint8(mode)
	mask := uint8((1 << bits) - 1)
	data &= mask
	mask <<= shift
	previous &= (^mask)
	previous |= data << shift
	// Enable temperature
	previous |= 0b10000000

	err := m.mmr.WriteUint8(m.datasheet.CRA_REG_M, previous)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)

	return nil
}

func (m *Magnetometer) GetRate() (MagnetometerRate, error) {
	value, err := m.mmr.ReadUint8(m.datasheet.CRA_REG_M)
	if err != nil {
		return MAGNETOMETER_RATE_30, err
	}
	const bits = 3
	const shift = 2
	range_ := ((uint32(value)) >> shift) & ((1 << bits) - 1)
	return MagnetometerRate(range_), nil
}

func (m *Magnetometer) SetGain(gain MagnetometerGain) error {
	const bits = 3
	const shift = 5

	data := uint8(gain)
	currentGain, err := m.mmr.ReadUint8(m.datasheet.CRB_REG_M)
	if err != nil {
		return err
	}

	mask := uint8((1 << bits) - 1)
	data &= mask
	mask <<= shift
	currentGain &= (^mask)
	currentGain |= data << shift
	err = m.mmr.WriteUint8(m.datasheet.CRB_REG_M, currentGain)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)

	m.gain = gain

	return nil
}

func (m *Magnetometer) GetGain() (MagnetometerGain, error) {
	value, err := m.mmr.ReadUint8(m.datasheet.CRA_REG_M)
	if err != nil {
		return MAGNETOMETER_GAIN_4_0, err
	}
	const bits = 3
	const shift = 5
	gain := ((uint32(value)) >> shift) & ((1 << bits) - 1)
	return MagnetometerGain(gain), nil
}

// The temperature sensor is technically on the same line as the magnetometer,
// so that's why I'm putting as a Magnetometer method. Note that the sensor is
// uncalibrated, so it can't return an absolute temperature, but from what I've
// read online, adding about 20 degrees C should get you close.
func (m *Magnetometer) SenseRelativeTemperature() (physic.Temperature, error) {
	degrees_eighths, err := m.senseRelativeTemperatureRaw()
	if err != nil {
		return 0, err
	}
	return physic.Temperature(int64(degrees_eighths)*int64(physic.Celsius)/8 + int64(physic.ZeroCelsius)), nil
}

// Returns the relative temperature in eights of a degree
func (m *Magnetometer) senseRelativeTemperatureRaw() (int16, error) {
	high, err := m.mmr.ReadUint8(m.datasheet.TEMP_OUT_H_M)
	if err != nil {
		return 0, err
	}
	low, err := m.mmr.ReadUint8(m.datasheet.TEMP_OUT_L_M)
	if err != nil {
		return 0, err
	}

	degreesEighths := ((int16(high) << 8) | int16(uint16(low))) >> 4
	return degreesEighths, nil
}
