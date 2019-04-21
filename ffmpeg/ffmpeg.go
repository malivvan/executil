package ffmpeg

import (
	"errors"
	"github.com/satnamram/executil"
	"strings"
)

func ToMp3(dst string, src string )error {
	out, err := executil.Command("ffmpeg", "-loglevel", "panic", "-i", src, dst).CombinedOutput()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		return errors.New(strings.Replace(string(out), "\n", " ", -1))
	}

	return nil
}