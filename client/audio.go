package client

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/audio"
)

const (
	sampleRate = 44100
	frequency  = 440
)

var (
	audioContext *audio.Context
	audioPlayer  *audio.Player
)

//infinite stream of a 440hz sine wave
type stream struct {
	position  int64
	remaining []byte
}

func init() {
	var err error
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}
	audioPlayer, err = audio.NewPlayer(audioContext, &stream{})
	if err != nil {
		log.Fatal(err)
	}
}

//io.Reader's Read - fills the data with sine wave samples
func (s *stream) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	const length = int64(sampleRate / frequency)
	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * max)
		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)
		p++
	}

	s.position += int64(len(buf))
	s.position %= length * 4

	if origBuf != nil {
		n := copy(origBuf, buf)
		s.remaining = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

//io.Closer's Close.
func (s *stream) Close() error {
	return nil
}

func Noise(on bool) {
	if on {
		if !audioPlayer.IsPlaying() {
			audioPlayer.Play()
		}
	} else {
		audioPlayer.Pause()
	}
}
