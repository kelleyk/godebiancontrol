// vim:ts=4:sw=4:noexpandtab
// © 2012-2014 Michael Stapelberg (see also: LICENSE)
package godebiancontrol_test

import (
	"bytes"
	"fmt"
	"github.com/stapelberg/godebiancontrol"
	"log"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	contents := bytes.NewBufferString(`Package: bti
Binary: bti
Version: 032-1
Maintainer: gregor herrmann <gregoa@debian.org>
Uploaders: tony mancill <tmancill@debian.org>
Build-Depends: debhelper (>= 8),
  bash-completion (>= 1:1.1-3),
  libcurl4-nss-dev, libreadline-dev, libxml2-dev, libpcre3-dev, liboauth-dev, xsltproc, docbook-xsl, docbook-xml, dh-autoreconf
Architecture: any
Standards-Version: 3.9.2
Format: 3.0 (quilt)
Files:
 3d5f65778bf3f89be03c313b0024b62c 1980 bti_032-1.dsc
 1e0d0b693fdeebec268004ba41701baf 59773 bti_032.orig.tar.gz
 ac1229a6d685023aeb8fcb0806324aa8 5065 bti_032-1.debian.tar.gz
Vcs-Browser: http://svn.toastfreeware.priv.at/wsvn/ToastfreewareDebian/bti/trunk/
Vcs-Svn: http://svn.toastfreeware.priv.at/debian/bti/trunk/
Checksums-Sha1:
 3da2c5a42138c884a7d9524b9706dc56c0d6d46e 1980 bti_032-1.dsc
 22061e3f56074703be415d65abc9ca27ef775c6a 59773 bti_032.orig.tar.gz
 66ae7f56a3c1f0ebe0638d0ec0599a819d72baea 5065 bti_032-1.debian.tar.gz
Checksums-Sha256:
 ed6015b79693f270d0a826c695b40e4d8eb4307942cac81a98f1fda479f74215 1980 bti_032-1.dsc
 feeabec98a89040a53283d798f7d55eb4311a854f17312a177dc45919883746a 59773 bti_032.orig.tar.gz
 f025da42efaf57db5e71a14cb8be27eb802ad23e7ab02b7ce2252454a86ac1d9 5065 bti_032-1.debian.tar.gz
Homepage: http://gregkh.github.com/bti/
Package-List:
 bti deb net extra
Directory: pool/main/b/bti
Priority: source
Section: net


Package: i3-wm
Version: 4.2-1
Installed-Size: 1573
Maintainer: Michael Stapelberg <stapelberg@debian.org>
Architecture: amd64
Provides: x-window-manager
Depends: libc6 (>= 2.8), libev4 (>= 1:4.04), libpcre3 (>= 8.10), libstartup-notification0 (>= 0.10), libx11-6, libxcb-icccm4 (>= 0.3.8), libxcb-keysyms1 (>= 0.3.8), libxcb-randr0 (>= 1.3), libxcb-util0 (>= 0.3.8), libxcb-xinerama0, libxcb1, libxcursor1 (>> 1.1.2), libyajl2 (>= 2.0.4), perl, x11-utils
Recommends: xfonts-base
Suggests: rxvt-unicode | x-terminal-emulator
Description: improved dynamic tiling window manager
 Key features of i3 are good documentation, reasonable defaults (changeable in
 a simple configuration file) and good multi-monitor support. The user
 interface is designed for power users and emphasizes keyboard usage. i3 uses
 XCB for asynchronous communication with X11 and aims to be fast and
 light-weight.
 .
 Please be aware i3 is primarily targeted at advanced users and developers.
Homepage: http://i3wm.org/
Description-md5: 2be7e62f455351435b1e055745d3e81c
Tag: implemented-in::c, interface::x11, role::program, uitoolkit::TODO,
 works-with::unicode, x11::window-manager
Section: x11
Priority: extra
Filename: pool/main/i/i3-wm/i3-wm_4.2-1_amd64.deb
Size: 798186
MD5sum: 3c7dbecd76d5c271401860967563fa8c
SHA1: 2e94f3faa5d4d617061f94076b2537d15fbff73f
SHA256: 2894bc999b3982c4e57f100fa31e21b52e14c5f3bc7ad5345f46842fcdab0db7`)
	paragraphs, err := godebiancontrol.Parse(contents)
	if err != nil {
		t.Fatal(err)
	}
	if len(paragraphs) != 2 {
		t.Fatal("Expected exactly two paragraphs")
	}
	if paragraphs[0]["Format"] != "3.0 (quilt)" {
		t.Fatal(`"Format" (simple) was not parsed correctly`)
	}
	if paragraphs[0]["Build-Depends"] != "debhelper (>= 8),bash-completion (>= 1:1.1-3),libcurl4-nss-dev, libreadline-dev, libxml2-dev, libpcre3-dev, liboauth-dev, xsltproc, docbook-xsl, docbook-xml, dh-autoreconf" {
		t.Fatal(`"Build-Depends" (folder) was not parsed correctly`)
	}

	expectedDescription := `improved dynamic tiling window manager
 Key features of i3 are good documentation, reasonable defaults (changeable in
 a simple configuration file) and good multi-monitor support. The user
 interface is designed for power users and emphasizes keyboard usage. i3 uses
 XCB for asynchronous communication with X11 and aims to be fast and
 light-weight.
 .
 Please be aware i3 is primarily targeted at advanced users and developers.`
	if paragraphs[1]["Description"] != expectedDescription {
		fmt.Print(expectedDescription)
		fmt.Print(paragraphs[1]["Description"])
		t.Fatal(`"Description" (multiline) was not parsed correctly`)
	}

	expectedFiles := `
 3d5f65778bf3f89be03c313b0024b62c 1980 bti_032-1.dsc
 1e0d0b693fdeebec268004ba41701baf 59773 bti_032.orig.tar.gz
 ac1229a6d685023aeb8fcb0806324aa8 5065 bti_032-1.debian.tar.gz`
	if paragraphs[0]["Files"] != expectedFiles {
		fmt.Print(expectedFiles)
		fmt.Print(paragraphs[0]["Files"])
		t.Fatal(`"Files" (multiline) was not parsed correctly`)
	}
}

func ExampleParse() {
	file, err := os.Open("debian-mirror/dists/sid/main/source/Sources")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	paragraphs, err := godebiancontrol.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	// Print a list of which source package uses which package format.
	for _, pkg := range paragraphs {
		fmt.Printf("%s uses %s\n", pkg["Package"], pkg["Format"])
	}
}
