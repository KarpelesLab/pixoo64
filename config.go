package pixoo64

// Config is the object returned by GetAllConf
type Config struct {
	Brightness          int
	RotationFlag        int // 1
	ClockTime           int // 60
	GalleryTime         int // 60
	SingleGalleyTime    int // -1
	PowerOnChannelId    int // 1
	GalleryShowTimeFlag int // 0|1
	CurClockId          int // 64
	Time24Flag          int
	TemperatureMode     int
	GyrateAngle         int
	MirrorFlag          int
	LightSwitch         int
}
