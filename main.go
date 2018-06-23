package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gjvnq/go-logger"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

var RootNode = &FNode{}
var FSConn *nodefs.FileSystemConnector
var FUSEServer *fuse.Server
var Unmounting bool
var Log *logger.Logger
var HideList []string
var SourcePath string

const DEFAULT_HIDE_LIST = "node_modules:.git:.cache:.svn:.hg"

func PrintCallDuration(prefix string, start *time.Time) {
	elapsed := time.Since(*start)
	Log.DebugNF(1, "%s: I took %s", prefix, elapsed)
}

func main() {
	var err error
	// Set Logger
	Log, err = logger.New("main", 1, os.Stdout)
	if err != nil {
		panic(err) // Check for error
	}

	// Get CLI options
	fuse_debug := flag.Bool("fuse-debug", false, "print debugging messages.")
	hide_list := flag.String("hide", DEFAULT_HIDE_LIST, "pattern for pretending files and folders don't exist.")
	other := flag.Bool("allow-other", false, "mount with -o allowother.")
	flag.Parse()
	if flag.NArg() < 2 {
		Log.FatalF("Usage:\n  filterfs SOURCE MOUNT_POINT")
	}
	SourcePath, _ = filepath.Abs(flag.Arg(0))
	mount_point, _ := filepath.Abs(flag.Arg(1))
	HideList = strings.Split(*hide_list, ":")
	RootNode.RealPath = SourcePath

	// Prepare fs
	FSConn = nodefs.NewFileSystemConnector(RootNode, &nodefs.Options{})
	mOpts := &fuse.MountOptions{
		AllowOther: *other,
		Name:       "filterfs",
		FsName:     mount_point,
		Debug:      *fuse_debug,
	}

	// Mount fs
	FUSEServer, err = fuse.NewServer(FSConn.RawFS(), mount_point, mOpts)
	if err != nil {
		Log.FatalF("Mount fail: %v", err)
	}

	// Prepare to deal with ctrl+c
	sig_chan := make(chan os.Signal, 20)
	signal.Notify(sig_chan, os.Interrupt)
	go func() {
		for _ = range sig_chan {
			Unmounting = true
			Log.Notice("Unmounting...")
			FUSEServer.Unmount()
			FUSEServer.Unmount()
			FUSEServer.Unmount()
			Log.Notice("Unmounted")
			os.Exit(0)
		}
	}()

	// Start things
	Log.Notice("Serving...")
	FUSEServer.Serve()
}
