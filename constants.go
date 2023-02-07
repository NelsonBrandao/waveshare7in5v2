package waveshare7in5v2

const (
	EPD_WIDTH  int = 800
	EPD_HEIGHT int = 480
)

const (
	PANEL_SETTING                byte = 0x00
	POWER_SETTING                byte = 0x01
	POWER_OFF                    byte = 0x02
	POWER_ON                     byte = 0x04
	BOOSTER_SOFT_START           byte = 0x06
	DEEP_SLEEP                   byte = 0x07
	DISPLAY_START_TRANSMISSION_1 byte = 0x10
	DISPLAY_REFRESH              byte = 0x12
	DISPLAY_START_TRANSMISSION_2 byte = 0x13
	DUAL_SPI                     byte = 0x15
	VCOM_DATA_INTERVAL_SETTING   byte = 0x50
	TCON                         byte = 0x60
	RESOLUTION_SETTING           byte = 0x61
	GET_STATUS                   byte = 0x71
	VCOM_DC                      byte = 0x82
)

const MAX_CHUNK_SIZE = 4096

const PIXEL_SIZE = 8
