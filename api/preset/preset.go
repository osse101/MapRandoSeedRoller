package preset

type Preset interface {
	Settings() ([]byte, error)
}
