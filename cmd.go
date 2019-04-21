package executil

import "os/exec"

func init() {

	// ffmpeg
	Register(Package{
		Name:      "ffmpeg",
		OS:        "linux",
		Arch:      "amd64",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz",
		Checksum:  "47b0d0cfd43e16f33e566832c1c78063d62e90aa217997573105eccc93327239",
		Installer: ArchiveInstaller("ffmpeg-4.1-64bit-static"),
		Activator: PathActivator,
	})

	// chromaprint
	Register(Package{
		Name:      "chromaprint",
		OS:        "linux",
		Arch:      "amd64",
		Commands:  []string{"fpcalc"},
		Download:  "https://github.com/acoustid/chromaprint/releases/download/v1.4.3/chromaprint-fpcalc-1.4.3-linux-x86_64.tar.gz",
		Checksum:  "a84425fccb43faa11b5bdc9d5b6101d6810b3b74876916191d42d31f7d73f4ce",
		Installer: ArchiveInstaller("chromaprint-fpcalc-1.4.3-linux-x86_64"),
		Activator: PathActivator,
	})
}


func Command(name string, arg ...string) *exec.Cmd {
	err := Ensure(name)
	if err != nil {
		info("warning:", err.Error())
	}
	return exec.Command("ffmpeg", arg...)
}