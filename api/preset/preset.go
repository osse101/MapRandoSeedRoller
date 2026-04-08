package preset

type Preset interface{
	Transform(input string)([]byte, error)
}
