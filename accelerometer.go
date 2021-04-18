package lsm303

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
	"periph.io/x/periph/conn/physic"
)

// This is a handle to the LSM303 accelerometer sensor.
type Accelerometer struct {
	mmr    mmr.Dev8
	sensorType SensorType
	datasheet *AccelerometerDatasheet
	addr *uint16
	range_ AccelerometerRange
	mode   AccelerometerMode
}

// New accelerometer opens a handle to an LSM303 accelerometer sensor.
func NewAccelerometer(bus i2c.Bus, opts ...AccelerometerOption) (*Accelerometer, error) {
	device := &Accelerometer{
		sensorType: LSM303DLHC,
		range_: ACCELEROMETER_RANGE_4G,
		mode:   ACCELEROMETER_MODE_NORMAL,
	}

	for i := range opts {
		opts[i].Apply(device)
	}

	if device.datasheet == nil {
		device.datasheet = datasheetForAccelerometer(device.sensorType)
	}

	if device.addr == nil {
		device.addr = &device.datasheet.ADDRESS
	}

	device.mmr = mmr.Dev8{
		Conn: &i2c.Dev{Bus: bus, Addr: uint16(*device.addr)},
		// I don't think we ever access more than 1 byte at once, so
		// this is irrelevant
		Order: binary.BigEndian,
	}

	// Enable the accelerometer 100 Hz, 0x57 = 0b01010111
	// Bits 0-2 = X, Y, Z enable
	// Bit 3 = low power mode
	// Bits 4-7 = speed, 0 = power down, 1-7 = 1 10 25 50 100 200 400 Hz, 8 = low
	//   power mode 1.62 khZ, 9 = normal 1.34 kHz / low power 5.376 kHz
	// TODO: Allow the user to set the Hz and toggle axes
	err := device.mmr.WriteUint8(device.datasheet.CTRL_REG1_A, 0x57)
	if err != nil {
		return nil, err
	}

	// Validate sensor
	if chipId, err := device.mmr.ReadUint8(device.datasheet.WHO_AM_I_A); err != nil {
		return nil, err
	} else if chipId != device.datasheet.CHIP_ID {
		return nil, fmt.Errorf("no %s detected", device.sensorType)
	}

	// Init accelerometer configuration
	device.SetRange(device.range_)
	device.SetMode(device.mode)

	return device, nil
}

func (a *Accelerometer) SenseRaw() (int16, int16, int16, error) {
	xLow, err := a.mmr.ReadUint8(a.datasheet.OUT_X_L_A)
	if err != nil {
		return 0, 0, 0, err
	}
	xHigh, err := a.mmr.ReadUint8(a.datasheet.OUT_X_H_A)
	if err != nil {
		return 0, 0, 0, err
	}
	yLow, err := a.mmr.ReadUint8(a.datasheet.OUT_Y_L_A)
	if err != nil {
		return 0, 0, 0, err
	}
	yHigh, err := a.mmr.ReadUint8(a.datasheet.OUT_Y_H_A)
	if err != nil {
		return 0, 0, 0, err
	}
	zLow, err := a.mmr.ReadUint8(a.datasheet.OUT_Z_L_A)
	if err != nil {
		return 0, 0, 0, err
	}
	zHigh, err := a.mmr.ReadUint8(a.datasheet.OUT_Z_H_A)
	if err != nil {
		return 0, 0, 0, err
	}

	xValue := int16(((uint16(xHigh)) << 8) + uint16(xLow))
	yValue := int16(((uint16(yHigh)) << 8) + uint16(yLow))
	zValue := int16(((uint16(zHigh)) << 8) + uint16(zLow))

	return xValue, yValue, zValue, nil
}

func (a *Accelerometer) Sense() (physic.Force, physic.Force, physic.Force, error) {
	xValue, yValue, zValue, err := a.SenseRaw()
	if err != nil {
		return 0, 0, 0, err
	}
	multiplier := getMultiplier(a.mode, a.range_)
	xAcceleration := (physic.Force)(int64(xValue) * multiplier)
	yAcceleration := (physic.Force)(int64(yValue) * multiplier)
	zAcceleration := (physic.Force)(int64(zValue) * multiplier)

	return xAcceleration, yAcceleration, zAcceleration, nil
}

func (a *Accelerometer) GetMode() (AccelerometerMode, error) {
	lowPowerU8, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG1_A)
	if err != nil {
		return ACCELEROMETER_MODE_NORMAL, err
	}
	lowPowerBit := readBits(uint32(lowPowerU8), 1, 3)

	highResolutionU8, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG4_A)
	if err != nil {
		return ACCELEROMETER_MODE_NORMAL, err
	}
	highResolutionBit := readBits(uint32(highResolutionU8), 1, 3)

	return AccelerometerMode((lowPowerBit << 1) | highResolutionBit), nil
}

func (a *Accelerometer) SetMode(mode AccelerometerMode) error {
	const bits = 1
	const shift = 3

	data := uint8((mode & 0x02) >> 1)
	power, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG1_A)
	if err != nil {
		return err
	}

	mask := uint8((1 << bits) - 1)
	data &= mask
	mask <<= shift
	power &= (^mask)
	power |= data << shift
	err = a.mmr.WriteUint8(a.datasheet.CTRL_REG1_A, power)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)

	data = uint8(mode & 0x01)
	resolution, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG4_A)
	if err != nil {
		return err
	}
	mask = uint8((1 << bits) - 1)
	data &= mask
	mask <<= shift
	resolution &= (^mask)
	resolution |= data << shift
	err = a.mmr.WriteUint8(a.datasheet.CTRL_REG4_A, resolution)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)

	a.mode = mode

	return nil
}

func (a *Accelerometer) GetRange() (AccelerometerRange, error) {
	value, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG4_A)
	if err != nil {
		return ACCELEROMETER_RANGE_4G, err
	}
	range_ := ((uint32(value)) >> 4) & ((1 << 2) - 1)
	return AccelerometerRange(range_), nil
}

func (a *Accelerometer) SetRange(range_ AccelerometerRange) error {
	const bits = 2
	const shift = 4

	data := uint8(range_)
	currentRange, err := a.mmr.ReadUint8(a.datasheet.CTRL_REG4_A)
	if err != nil {
		return err
	}

	mask := uint8((1 << bits) - 1)
	data &= mask
	mask <<= shift
	currentRange &= (^mask)
	currentRange |= data << shift
	err = a.mmr.WriteUint8(a.datasheet.CTRL_REG4_A, currentRange)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 20)

	a.range_ = range_

	return nil
}

func (a *Accelerometer) String() string {
	return "LSM303 accelerometer"
}

func readBits(value uint32, bits uint32, shift uint8) uint32 {
	value >>= shift
	return value & ((1 << bits) - 1)
}

// Gets the multiplier for the accelerometer mode and range
func getMultiplier(mode AccelerometerMode, range_ AccelerometerRange) int64 {
	// The constants in here needed to be rounded because some of then aren't
	// exactly representable. I added tests for what the true value should be.
	switch mode {
	case ACCELEROMETER_MODE_LOW_POWER:
		switch range_ {
		case ACCELEROMETER_RANGE_2G:
			return 153277939 >> 8
		case ACCELEROMETER_RANGE_4G:
			return 306555879 >> 8
		case ACCELEROMETER_RANGE_8G:
			return 613111758 >> 8
		case ACCELEROMETER_RANGE_16G:
			return 1839531407 >> 8
		}
	case ACCELEROMETER_MODE_NORMAL:
		switch range_ {
		case ACCELEROMETER_RANGE_2G:
			return 38245935 >> 6
		case ACCELEROMETER_RANGE_4G:
			return 76688003 >> 6
		case ACCELEROMETER_RANGE_8G:
			return 153277939 >> 6
		case ACCELEROMETER_RANGE_16G:
			return 459931885 >> 6
		}

	case ACCELEROMETER_MODE_HIGH_RESOLUTION:
		switch range_ {
		case ACCELEROMETER_RANGE_2G:
			return 9610517 >> 4
		case ACCELEROMETER_RANGE_4G:
			return 19122967 >> 4
		case ACCELEROMETER_RANGE_8G:
			return 38245935 >> 4
		case ACCELEROMETER_RANGE_16G:
			return 114933938 >> 4
		}
	default:
		log.Fatalf("Unknown mode %v in getMultiplier", mode)
	}
	log.Fatalf("Unknown range %v in getMultiplier", range_)
	return 0.0
}
