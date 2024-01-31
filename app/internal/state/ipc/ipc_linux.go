//go:build (linux && !android) || freebsd || openbsd

package ipc

import (
	"errors"
	"fmt"
	"github.com/partisiadev/partisiawallet/app/assets"

	log "github.com/partisiadev/partisiawallet/log"
	syscall "golang.org/x/sys/unix"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

/* Format of appName.desktop entry file
[Desktop Entry]
Version=Application's versionInEntryFile default is 1.0.0
Type=Application
Name=Your App Name
Exec=Path to executable/binary file %U
Icon=Path to iconPath
MimeType=comma separated schemes (mimeType)
StartupNotify=bool(currently always true)
Terminal=bool(currently always false)
SingleMainWindow=bool(currently always true)"
*/

/* Default paths used
~/.local/share
~/.local/share/applications
~/.local/share/applications/entryFile.desktop
~/.local/share/appDataDir
~/.local/share/appDataDir/icons (app icons dir path)
~/.local/share/appDataDir/bin (app binaries dir path)
/tmp/socket
*/

const DefaultMimeType = "x-scheme-handler/partisiawallet"
const DefaultAppBinaryName = "partisiawallet"
const DefaultAppNameInEntryFile = "Partisia Wallet"
const DefaultSocketDirPath = "/tmp"
const DefaultSocketFileName = "partisiawallet.sock"
const DefaultDesktopEntryFileName = "Partisiawallet.desktop"
const DefaultVersionInEntryFile = "0.0.0"

// mimeType is comma separated schemes,this is the
// only value required for deep linking
// ex -ldflags="-X 'gioui.org/app.mimeType=x-scheme-handler/custom-uri'"

// Default is DefaultMimeType
var mimeType string
var socketConn net.Listener = nil

// socketDirPath Default is DefaultSocketDirPath
var socketDirPath string

// socketFileName Default is DefaultSocketFileName
var socketFileName string

// desktopEntryDir
// Default is "/user'sHomeDir/.local/share/applications"
var desktopEntryDir string

// desktopEntryFileName
// name of the entry file
// Default is DefaultDesktopEntryFileName
var desktopEntryFileName string

// appDataDir
// Default is "user'sHomeDir/.local/share/partisiawallet"
var appDataDir string

// binDirPath
// Default is "/user'sHomeDir/.local/share/partisiawallet/bin"
var binDirPath string

// appBinName
// Default is DefaultAppBinaryName
var appBinName string

// iconsDirPath
// Default is "/user'sHomeDir/.local/share/partisiawallet/icons"
var iconsDirPath string

// versionInEntryFile for desktop entry file
// Default is DefaultVersionInEntryFile
var versionInEntryFile string

// nameInEntryFile for desktop entry file
// Default is DefaultAppNameInEntryFile
var nameInEntryFile string

// FullName implies full path to socket file
var socketFileFullName string

// iconPath
// Icon from iconPath is copied to iconsDirPath and
// new icon path is added to desktop entry file.
var iconPath string

func init() {
	if mimeType == "" {
		mimeType = DefaultMimeType
	}
	if appBinName == "" {
		appBinName = DefaultAppBinaryName
	}
	if socketDirPath == "" {
		socketDirPath = DefaultSocketDirPath
	}
	if socketFileName == "" {
		socketFileName = DefaultSocketFileName
	}

	socketFileFullName = path.Join(socketDirPath, socketFileName)
	if !strings.HasSuffix(socketFileFullName, ".sock") {
		socketFileFullName += ".sock"
	}
	c, err := net.Dial("unix", socketFileFullName)
	if err != nil {
		// syscall.ECONNREFUSED error most likely indicates socket file exists but
		//  app instance is not running, hence we delete the socketFile
		if errors.Is(err, syscall.ECONNREFUSED) {
			// delete socket file
			_ = os.Remove(socketFileFullName)
		}
		// we exit with error if error is other than these errors
		// (syscall.ENOENT indicates that socket file doesn't exist)
		if !errors.Is(err, syscall.ECONNREFUSED) && !errors.Is(err, syscall.ENOENT) {
			log.Logger().Fatal(err)
		}
		err = nil
	} else {
		// since err is nil, we are certain that another instance of our app is running
		// if any arguments were passed to this app, then we pass it to already running
		// instance of our app
		if len(os.Args) > 1 {
			_, _ = c.Write([]byte(strings.Join(os.Args[:], "\n")))
		}
		_ = c.Close()
		log.Logger().Fatal("another instance of app is already running")
	}
	socketConn, err = net.Listen("unix", socketFileFullName)
	if err != nil {
		log.Logger().Fatal(err)
	}
	err = createDesktopEntry()
	if err != nil {
		_ = socketConn.Close()
		log.Logger().Fatal(err)
	}
	go listenToSocketConn()
}

// listenToSocketConn blocking function
func listenToSocketConn() {
	if socketConn == nil {
		return
	}
	// Cleanup the socketFile.
	defer func() {
		_ = socketConn.Close()
		_ = os.Remove(socketFileFullName)
	}()
	for {
		// Accept an incoming connection.
		conn, err := socketConn.Accept()
		if err != nil {
			continue
		}
		log.Logger().Println(conn.LocalAddr())
		log.Logger().Println(conn.RemoteAddr())

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer func() {
				_ = conn.Close()
			}()
			bs, err := io.ReadAll(conn)
			if err != nil {
				return
			}
			args := strings.Split(string(bs), "\n")
			log.Logger().Printf("%#v\n", args)
		}(conn)
	}
}

func createDesktopEntry() (err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Logger().Fatal(err)
		return err
	}
	dataHome := filepath.Join(homeDir, ".local", "share")
	if desktopEntryDir == "" {
		desktopEntryDir = filepath.Join(dataHome, "applications")
	}
	if appDataDir == "" {
		appDataDir = filepath.Join(dataHome, appBinName)
	}
	if binDirPath == "" {
		binDirPath = filepath.Join(appDataDir, "bin")
	}
	if iconsDirPath == "" {
		iconsDirPath = filepath.Join(appDataDir, "icons")
	}
	if versionInEntryFile == "" {
		versionInEntryFile = DefaultVersionInEntryFile
	}
	if nameInEntryFile == "" {
		nameInEntryFile = DefaultAppNameInEntryFile
	}
	// Create applications directory if not exists.
	if _, err = os.Stat(desktopEntryDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(desktopEntryDir, 0755)
		if err != nil {
			log.Logger().Fatal(err)
			return
		}
	}
	// Create bin directory if not exists
	if _, err = os.Stat(binDirPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(binDirPath, 0755)
		if err != nil {
			log.Logger().Fatal(err)
			return
		}
	}
	// Create icons directory if not exists
	if _, err = os.Stat(iconsDirPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(iconsDirPath, 0755)
		if err != nil {
			log.Logger().Fatal(err)
			return
		}
	}
	// copy iconPath from iconPath to icon entry in entry file
	var src, dst *os.File
	var srcAlt fs.File
	var iconPathEmpty bool
	if iconPath != "" {
		defer func(src *os.File) {
			err := src.Close()
			if err != nil {
				log.Logger().Errorln(err)
			}
		}(src)
		src, err = os.Open(iconPath)
		if err != nil {
			log.Logger().Fatal(err)
			return err
		}
		iconPath = filepath.Join(iconsDirPath, filepath.Base(iconPath))
	} else {
		defer func() {
			err := srcAlt.Close()
			if err != nil {
				log.Logger().Errorln(err)
			}
		}()
		srcAlt, err = assets.AppIconFile.Open(assets.AppIconFileName)
		if err != nil {
			log.Logger().Fatal(err)
		}
		iconPathEmpty = true
		iconPath = filepath.Join(iconsDirPath, assets.AppIconFileName)
	}
	dst, err = os.Create(iconPath)
	if err != nil {
		log.Logger().Fatal(err)
		return err
	}
	defer func() {
		err := dst.Close()
		if err != nil {
			log.Logger().Errorln(err)
		}
	}()
	if !iconPathEmpty {
		_, err = io.Copy(dst, src)
	} else {
		_, err = io.Copy(dst, srcAlt)
	}
	if err != nil {
		log.Logger().Fatal(err)
		return err
	}
	binFilePath := filepath.Join(binDirPath, appBinName)
	// copy only if the src and dest binaries path are different
	if binFilePath != os.Args[0] {
		binFileBs, err := os.ReadFile(os.Args[0])
		if err != nil {
			log.Logger().Fatal(err)
			return err
		}
		err = os.WriteFile(binFilePath, binFileBs, 0755)
		// only if error is not 'text file busy'
		if err != nil && !errors.Is(err, syscall.ETXTBSY) {
			log.Logger().Fatal(err)
			return err
		}
	}
	entryFile := fmt.Sprintf(
		"[Desktop Entry]\n"+
			"Version=%s\n"+
			"Type=Application\n"+
			"Name=%s\n"+
			"Exec=%s %%U\n"+
			"Icon=%s\n"+
			"MimeType=%s\n"+
			"StartupNotify=true\n"+
			"Terminal=false\n"+
			"SingleMainWindow=true\n",
		versionInEntryFile,
		nameInEntryFile,
		binFilePath,
		iconPath,
		DefaultMimeType,
	)
	if desktopEntryFileName == "" {
		desktopEntryFileName = DefaultDesktopEntryFileName
	}
	if !strings.HasSuffix(desktopEntryFileName, ".desktop") {
		desktopEntryFileName += ".desktop"
	}
	desktopEntryFilePath := filepath.Join(desktopEntryDir, desktopEntryFileName)
	err = os.WriteFile(desktopEntryFilePath, []byte(entryFile), 0755)
	if err != nil {
		log.Logger().Fatal(err)
		return err
	}
	err = exec.Command("update-desktop-database", desktopEntryDir).Start()
	if err != nil {
		log.Logger().Fatal(err)
		return
	}
	return nil
}
