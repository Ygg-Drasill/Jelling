package ui

type JellTheme struct {
	Primary     string `mapstructure:"primary"`
	Secondary   string `mapstructure:"secondary"`
	Text        string `mapstructure:"text"`
	Placeholder string `mapstructure:"placeholder"`
}

var Theme JellTheme = JellTheme{
	Primary:     "179",
	Secondary:   "117",
	Text:        "230",
	Placeholder: "59",
}
