package lsm303

type AccelerometerDatasheet struct {
	ADDRESS     uint16
	WHO_AM_I_A  uint8
	CHIP_ID     uint8
	CTRL_REG1_A uint8
	// CTRL_REG2_A uint8
	// CTRL_REG3_A uint8
	CTRL_REG4_A uint8
	// CTRL_REG5_A uint8
	// CTRL_REG6_A uint8
	// REFERENCE_A uint8
	// STATUS_REG_A  uint8
	OUT_X_L_A uint8
	OUT_X_H_A uint8
	OUT_Y_L_A uint8
	OUT_Y_H_A uint8
	OUT_Z_L_A uint8
	OUT_Z_H_A uint8
	// FIFO_CTRL_REG_A uint8
	// FIFO_SRC_REG_A  uint8
	// INT1_CFG_A      uint8
	// INT1_SOURCE_A   uint8
	// INT1_THS_A      uint8
	// INT1_DURATION_A uint8
	// INT2_CFG_A      uint8
	// INT2_SOURCE_A   uint8
	// INT2_THS_A      uint8
	// INT2_DURATION_A uint8
	// CLICK_CFG_A     uint8
	// CLICK_SRC_A     uint8
	// CLICK_THS_A     uint8
	// TIME_LATENCY_A  uint8
	// TIME_LIMIT_A    uint8
	// TIME_WINDOW_A   uint8
}

type MagnetometerDatasheet struct {
	ADDRESS    uint16
	WHO_AM_I_M uint8
	CHIP_ID    uint8
	CRA_REG_M  uint8
	CRB_REG_M  uint8
	MR_REG_M   uint8
	OUT_X_H_M  uint8
	OUT_X_L_M  uint8
	OUT_Z_H_M  uint8
	OUT_Z_L_M  uint8
	OUT_Y_H_M  uint8
	OUT_Y_L_M  uint8
	// SR_REG_M uint8
	IRA_REG_M uint8
	// IRB_REG_M uint8
	// IRC_REG_M uint8
	TEMP_OUT_H_M uint8
	TEMP_OUT_L_M uint8
}

func datasheetForAccelerometer(sensorType SensorType) *AccelerometerDatasheet {
	datasheet := &AccelerometerDatasheet{
		ADDRESS:     0x19,
		WHO_AM_I_A:  0x0F,
		CHIP_ID:     0x33,
		CTRL_REG1_A: 0x20,
		// CTRL_REG2_A:     0x21,
		// CTRL_REG3_A:     0x22,
		CTRL_REG4_A: 0x23,
		// CTRL_REG5_A:     0x24,
		// CTRL_REG6_A:     0x25,
		// REFERENCE_A:     0x26,
		// STATUS_REG_A:    0x27,
		OUT_X_L_A: 0x28,
		OUT_X_H_A: 0x29,
		OUT_Y_L_A: 0x2A,
		OUT_Y_H_A: 0x2B,
		OUT_Z_L_A: 0x2C,
		OUT_Z_H_A: 0x2D,
		// FIFO_CTRL_REG_A: 0x2E,
		// FIFO_SRC_REG_A:  0x2F,
		// INT1_CFG_A:      0x30,
		// INT1_SOURCE_A:   0x31,
		// INT1_THS_A:      0x32,
		// INT1_DURATION_A: 0x33,
		// INT2_CFG_A:      0x34,
		// INT2_SOURCE_A:   0x35,
		// INT2_THS_A:      0x36,
		// INT2_DURATION_A: 0x37,
		// CLICK_CFG_A:     0x38
		// CLICK_SRC_A:     0x39,
		// CLICK_THS_A:     0x3A,
		// TIME_LATENCY_A:  0x3B,
		// TIME_LIMIT_A:    0x3C
		// TIME_WINDOW_A:   0x3D,
	}
	switch sensorType {
	case LSM303DLHC:
		return datasheet
	case LSM303AGR:
		return datasheet
	case LSM303C:
		datasheet.ADDRESS = 0x1D
		datasheet.CHIP_ID = 0x41
		return datasheet
	default:
		return datasheet
	}
}

func datasheetForMagnetometer(sensorType SensorType) *MagnetometerDatasheet {
	defaultDatasheet := &MagnetometerDatasheet{
		ADDRESS: 0x1E,
		// The LSM303(DLHC) magnetometer doesn't have an ID register,
		// IRA_REG_M used instead with constant value: 0b01001000
		WHO_AM_I_M: 0x0A,
		CHIP_ID:    0b01001000,
		CRA_REG_M:  0x00,
		CRB_REG_M:  0x01,
		MR_REG_M:   0x02,
		OUT_X_H_M:  0x03,
		OUT_X_L_M:  0x04,
		OUT_Z_H_M:  0x05,
		OUT_Z_L_M:  0x06,
		OUT_Y_H_M:  0x07,
		OUT_Y_L_M:  0x08,
		// SR_REG_M:     0x09,
		IRA_REG_M: 0x0A,
		// IRB_REG_M:    0x0B,
		// IRC_REG_M:    0x0C,
		TEMP_OUT_H_M: 0x31,
		TEMP_OUT_L_M: 0x32,
	}

	switch sensorType {
	case LSM303DLHC:
		return defaultDatasheet
	case LSM303AGR:
		return &MagnetometerDatasheet{
			ADDRESS:    0x1E,
			WHO_AM_I_M: 0x4F,
			CHIP_ID:    0x40,
			CRA_REG_M:  0x60,
			CRB_REG_M:  0x61,
			MR_REG_M:   0x02,
			OUT_X_L_M:  0x68,
			OUT_X_H_M:  0x69,
			OUT_Y_L_M:  0x6A,
			OUT_Y_H_M:  0x6B,
			OUT_Z_L_M:  0x6C,
			OUT_Z_H_M:  0x6D,
			IRA_REG_M:  0x0A,
			TEMP_OUT_H_M: 0x31, // Couldn't verify if this sensor is able to measure temperature
			TEMP_OUT_L_M: 0x32,
		}
	case LSM303C:
		return &MagnetometerDatasheet{
			ADDRESS:    0x1E,
			WHO_AM_I_M: 0x0F,
			CHIP_ID:    0x3D,
			CRA_REG_M:  0x20, // Called CTRL_REG1_M in LSM303C Datasheet
			MR_REG_M:   0x22, // Called CTRL_REG3_M in LSM303C Datasheet
			OUT_X_L_M:  0x28,
			OUT_X_H_M:  0x29,
			OUT_Y_L_M:  0x2A,
			OUT_Y_H_M:  0x2B,
			OUT_Z_L_M:  0x2C,
			OUT_Z_H_M:  0x2D,
			IRA_REG_M:  0x0A,
			TEMP_OUT_H_M: 0x2F,
			TEMP_OUT_L_M: 0x2E,
		}
	default:
		return defaultDatasheet
	}
}
