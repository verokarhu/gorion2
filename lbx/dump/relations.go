package lbx

type imagePair struct {
	Filename string
	Palette  string
}

type palettePair struct {
	Filename string
	Index    int
}

var video_files = []string{
	"amebafin",
	"anatkfin",
	"anwinfin",
	"dimtvfin",
	"genwinfn",
	"intro",
	"loserfin",
	"orionfin",
	"plntdfin",
	"wininfin",
}

var audio_files = []string{
	"sound",
	"stream",
	"streamhd",
}

var palette_files = []palettePair{
	{"fonts", 1},
	{"fonts", 2},
	{"fonts", 3},
	{"fonts", 4},
	{"fonts", 5},
	{"fonts", 6},
	{"fonts", 7},
	{"fonts", 8},
	{"fonts", 9},
	{"fonts", 10},
	{"fonts", 11},
	{"fonts", 12},
	{"fonts", 13},
	{"ifonts", 1},
	{"ifonts", 2},
	{"ifonts", 3},
	{"ifonts", 4},
}

var image_files = []imagePair{
	{"mainmenu", "fonts6"},
	//	{"planets", "none"},
	//	{"logo", "none"},
	//	{"bldg0", "all"},
	//	{"app_pics", "none"},
	//	{"newgame", "none"},
	//	{"racesel", "none"},
	//	{"tanm_001", "all"},
}
