package workflow

func GetHelp(source string) string {
	helpText := "!roll <preset> <flags>. This season's presets are s5, expert, and mentor.  Add starting items by adding 'S:ITEM1ITEM2' in caps."
	//TODO: Add source handling
	return helpText
}
