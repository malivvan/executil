package executil

import "os/exec"

func init() {

	// ffmpeg static for linux from www.johnvansickle.com
	// (md5sum was checked before creating sha256sum)
	Register(Package{
		Name:      "ffmpeg",
		OS:        "linux",
		Arch:      "amd64",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://www.johnvansickle.com/ffmpeg/old-releases/ffmpeg-4.0.3-64bit-static.tar.xz",
		Checksum:  "0877b4945e0963e4b0607858faa6b320940cc5f1e06253e18dce950511d63ac3",
		Installer: ArchiveInstaller("ffmpeg-4.0.3-64bit-static"),
		Activator: PathActivator,
	})
	Register(Package{
		Name:      "ffmpeg",
		OS:        "linux",
		Arch:      "386",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://www.johnvansickle.com/ffmpeg/old-releases/ffmpeg-4.0.3-32bit-static.tar.xz",
		Checksum:  "bcb679db85c574314e2431a63f6430eed23cc8e38f74806f2f1d71e7cd16cb34",
		Installer: ArchiveInstaller("ffmpeg-4.0.3-32bit-static"),
		Activator: PathActivator,
	})

	// ffmpeg static for windows from ffmpeg.zeranoe.com
	// (no checksums provided, sha256sum was created after downloading via https)
	Register(Package{
		Name:      "ffmpeg",
		OS:        "windows",
		Arch:      "386",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://ffmpeg.zeranoe.com/builds/win32/static/ffmpeg-4.1.1-win32-static.zip",
		Checksum:  "b739e7d1eff03f4215858aad3a10393f12d2f8843831870e7c148d04692e012c",
		Installer: ArchiveInstaller("ffmpeg-4.1.1-win32-static/bin"),
		Activator: PathActivator,
	})
	Register(Package{
		Name:      "ffmpeg",
		OS:        "windows",
		Arch:      "amd64",
		Commands:  []string{"ffmpeg", "ffprobe"},
		Download:  "https://ffmpeg.zeranoe.com/builds/win64/static/ffmpeg-4.1.1-win64-static.zip",
		Checksum:  "2c19658d69de08ea7ef585fbf801ef6a364795b5288e23d79066876ff0465df6",
		Installer: ArchiveInstaller("ffmpeg-4.1.1-win64-static/bin"),
		Activator: PathActivator,
	})

	// chromaprint from github.com/acoustid/chromaprint/releases/
	// (no checksums provided, sha256sum was created after downloading via https)
	Register(Package{
		Name:      "chromaprint",
		OS:        "linux",
		Arch:      "386",
		Commands:  []string{"fpcalc"},
		Download:  "https://github.com/acoustid/chromaprint/releases/download/v1.4.3/chromaprint-fpcalc-1.4.3-linux-i686.tar.gz",
		Checksum:  "96fec7e564b46f373c3ddffd52ed76db70a1f34f5331bdd4e8ad288db4b60ee8",
		Installer: ArchiveInstaller("chromaprint-fpcalc-1.4.3-linux-i686"),
		Activator: PathActivator,
	})
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
	Register(Package{
		Name:      "chromaprint",
		OS:        "windows",
		Arch:      "386",
		Commands:  []string{"fpcalc"},
		Download:  "https://github.com/acoustid/chromaprint/releases/download/v1.4.3/chromaprint-fpcalc-1.4.3-windows-i686.zip",
		Checksum:  "cd9e45581fd075f1e85eb55966bbf251115aade289191ad9131e1463597e6e98",
		Installer: ArchiveInstaller("chromaprint-fpcalc-1.4.3-windows-i686"),
		Activator: PathActivator,
	})
	Register(Package{
		Name:      "chromaprint",
		OS:        "windows",
		Arch:      "amd64",
		Commands:  []string{"fpcalc"},
		Download:  "https://github.com/acoustid/chromaprint/releases/download/v1.4.3/chromaprint-fpcalc-1.4.3-windows-x86_64.zip",
		Checksum:  "7d904a95c0d8738973c6bde55968f53dcc95c72597431e11619d355c33edc199",
		Installer: ArchiveInstaller("chromaprint-fpcalc-1.4.3-windows-x86_64"),
		Activator: PathActivator,
	})
}

func Command(name string, arg ...string) *exec.Cmd {
	err := Ensure(name)
	if err != nil {
		info("warning:", err.Error())
	}
	return exec.Command(name, arg...)
}


