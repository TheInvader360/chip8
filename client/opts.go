package client

type Opts struct {
	Path  string //path to rom file
	Clock int    //cpu clock speed in hz
	ScrW  int    //width of client screen in pixels
	ScrH  int    //height of client screen in pixels
	ScrF  bool   //is fullscreen mode enabled
	Debug bool   //is debug enabled
}
