package ffmpeg

import (
	"errors"
	"github.com/satnamram/shtool"
	"os/exec"
	"strings"
)

func init() {
	shtool.Register(shtool.Package{
		Name:      "ffmpeg",
		OS:        "linux",
		Arch:      "amd64",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz",
		Checksum:  "47b0d0cfd43e16f33e566832c1c78063d62e90aa217997573105eccc93327239",
		Installer: shtool.ArchiveInstaller("ffmpeg-4.1-64bit-static"),
		Activator: shtool.PathActivator,
	})
}

func ToMp3(dst string, src string )error {
	err := shtool.Ensure("ffmpeg")
	if err != nil {
		return err
	}

	out, err := exec.Command("ffmpeg", "-loglevel", "panic", "-i", src, dst).CombinedOutput()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		return errors.New(strings.Replace(string(out), "\n", " ", -1))
	}

	return nil
}