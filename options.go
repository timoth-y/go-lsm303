package lsm303


type (
	// AccelerometerOption configures a LSM303 accelerometer.
	AccelerometerOption interface {
		Apply(*Accelerometer)
	}
	// AccelerometerOptionFunc is a function that configures a device.
	AccelerometerOptionFunc func(*Accelerometer)

	// MagnetometerOption configures a LSM303 magnetometer.
	MagnetometerOption interface {
		Apply(*Magnetometer)
	}
	// MagnetometerOptionFunc is a function that configures a device.
	MagnetometerOptionFunc func(*Magnetometer)
)

type SensorType string

const (
	LSM303DLHC SensorType = "LSM303DLHC"
	LSM303AGR SensorType = "LSM303AGR"
	LSM303C SensorType = "LSM303C"
)

type AccelerometerMode int

const (
	ACCELEROMETER_MODE_NORMAL AccelerometerMode = iota
	ACCELEROMETER_MODE_HIGH_RESOLUTION
	ACCELEROMETER_MODE_LOW_POWER
)

type AccelerometerRange int

const (
	ACCELEROMETER_RANGE_2G AccelerometerRange = iota
	ACCELEROMETER_RANGE_4G
	ACCELEROMETER_RANGE_8G
	ACCELEROMETER_RANGE_16G
)

func (mode AccelerometerMode) String() string {
	return [...]string{"normal", "high resolution", "low power"}[mode]
}

func (range_ AccelerometerRange) String() string {
	return [...]string{"2G", "4G", "8G", "16G"}[range_]
}


// Apply calls OptionFunc on device instance
func (f AccelerometerOptionFunc) Apply(dev *Accelerometer) {
	f(dev)
}

// Apply calls OptionFunc on device instance
func (f MagnetometerOptionFunc) Apply(dev *Magnetometer) {
	f(dev)
}

// WithAccelerometerSensorType can be used to specify LSM303 family sensor type.
// Default is LSM303DLHC.
func WithAccelerometerSensorType(sensorType SensorType) AccelerometerOption {
	return AccelerometerOptionFunc(func(d *Accelerometer) {
		d.sensorType = sensorType
	})
}

// WithAccelerometerAddress can be used to specify I²C address for Accelerometer.
// Default is 0x19 for LSM303 and 0x1E for LSM303C.
func WithAccelerometerAddress(addr uint16) AccelerometerOption {
	return AccelerometerOptionFunc(func(d *Accelerometer) {
		d.addr = &addr
	})
}

// WithMode can be used to specify accelerometer mode.
// Default is ACCELEROMETER_MODE_NORMAL.
func WithMode(mode AccelerometerMode) AccelerometerOption {
	return AccelerometerOptionFunc(func(d *Accelerometer) {
		d.mode = mode
	})
}

// WithRange can be used to specify accelerometer range.
// Default is ACCELEROMETER_RANGE_4G.
func WithRange(rng AccelerometerRange) AccelerometerOption {
	return AccelerometerOptionFunc(func(d *Accelerometer) {
		d.range_ = rng
	})
}

type MagnetometerGain int

const (
	MAGNETOMETER_GAIN_1_3 MagnetometerGain = iota
	MAGNETOMETER_GAIN_1_9
	MAGNETOMETER_GAIN_2_5
	MAGNETOMETER_GAIN_4_0
	MAGNETOMETER_GAIN_4_7
	MAGNETOMETER_GAIN_5_6
	MAGNETOMETER_GAIN_8_1
)

func (mode MagnetometerGain) String() string {
	return [...]string{"1.3", "1.9", "2.5", "4.0", "4.7", "5.6", "8.1"}[mode]
}

type MagnetometerRate int

const (
	MAGNETOMETER_RATE_0_75 MagnetometerRate = iota
	MAGNETOMETER_RATE_1_5
	MAGNETOMETER_RATE_3_0
	MAGNETOMETER_RATE_7_5
	MAGNETOMETER_RATE_15
	MAGNETOMETER_RATE_30
	MAGNETOMETER_RATE_75
	MAGNETOMETER_RATE_220
)

func (range_ MagnetometerRate) String() string {
	return [...]string{"0.75", "1.55", "3.05", "7.55", "15", "30", "75", "220"}[range_]
}

// WithMagnetometerSensorType can be used to specify LSM303 family sensor type.
// Default is LSM303DLHC.
func WithMagnetometerSensorType(sensorType SensorType) MagnetometerOption {
	return MagnetometerOptionFunc(func(d *Magnetometer) {
		d.sensorType = sensorType
	})
}

// WithMagnetometerAddress can be used to specify I²C address for Magnetometer.
// Default is 0x1E.
func WithMagnetometerAddress(addr uint16) MagnetometerOption {
	return MagnetometerOptionFunc(func(d *Magnetometer) {
		d.addr = &addr
	})
}

// WithGain can be used to specify magnetometer gain.
// Default is MAGNETOMETER_GAIN_4_0.
func WithGain(gain MagnetometerGain) MagnetometerOption {
	return MagnetometerOptionFunc(func(d *Magnetometer) {
		d.gain = gain
	})
}

// WithRange can be used to specify magnetometer rate.
// Default is MAGNETOMETER_RATE_30.
func WithRate(rate MagnetometerRate) MagnetometerOption {
	return MagnetometerOptionFunc(func(d *Magnetometer) {
		d.rate = rate
	})
}

// WithDatasheet can be used to specify datasheet addresses,
// in case new LSM family device appears.
func WithDatasheet(datasheet MagnetometerDatasheet) MagnetometerOption {
	return MagnetometerOptionFunc(func(d *Magnetometer) {
		d.datasheet = &datasheet
	})
}
