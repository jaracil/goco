// Package file is a GopherJS wrapper for cordova file plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-file
package file

import (
	"errors"
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

var (
	ErrNotFound              = errors.New("Not found error")
	ErrSecurity              = errors.New("Security error")
	ErrAbort                 = errors.New("Abort error")
	ErrNotReadable           = errors.New("Not readable error")
	ErrEncoding              = errors.New("Encoding error")
	ErrNoModificationAllowed = errors.New("No modification allowed error")
	ErrInvalidState          = errors.New("Invalid state error")
	ErrSyntax                = errors.New("Syntax error")
	ErrInvalidModification   = errors.New("Invalid modification error")
	ErrQuotaExceeded         = errors.New("Quota exceeded error")
	ErrTypeMismatch          = errors.New("Type mismatch error")
	ErrPathExists            = errors.New("Path exists error")
)

const (
	NotFoundErrVal              = 1
	SecurityErrVal              = 2
	AbortErrVal                 = 3
	NotReadableErrVal           = 4
	EncodingErrVal              = 5
	NoModificationAllowedErrVal = 6
	InvalidStateErrVal          = 7
	SyntaxErrVal                = 8
	InvalidModificationErrVal   = 9
	QuotaExceededErrVal         = 10
	TypeMismatchErrVal          = 11
	PathExistsErrVal            = 12
)

func (fe *FileError) Error() error {
	switch fe.Code {
	case NotFoundErrVal:
		return ErrNotFound
	case SecurityErrVal:
		return ErrSecurity
	case AbortErrVal:
		return ErrAbort
	case NotReadableErrVal:
		return ErrNotReadable
	case EncodingErrVal:
		return ErrEncoding
	case NoModificationAllowedErrVal:
		return ErrNoModificationAllowed
	case InvalidStateErrVal:
		return ErrInvalidState
	case SyntaxErrVal:
		return ErrSyntax
	case InvalidModificationErrVal:
		return ErrInvalidModification
	case QuotaExceededErrVal:
		return ErrQuotaExceeded
	case TypeMismatchErrVal:
		return ErrTypeMismatch
	case PathExistsErrVal:
		return ErrPathExists
	}
	return fmt.Errorf("Unknown error: %d", fe.Code)
}

type Directories struct {
	*js.Object

	// ApplicationDirectory holds read-only directory where the application is installed. (iOS, Android, BlackBerry 10, OSX, windows)
	ApplicationDirectory string `js:"applicationDirectory"`

	// ApplicationStorageDirectory holds root directory of the application's sandbox;
	// on iOS & windows this location is read-only (but specific subdirectories [like /Documents on iOS or /localState on windows] are read-write).
	// All data contained within is private to the app. (iOS, Android, BlackBerry 10, OSX)
	ApplicationStorageDirectory string `js:"applicationStorageDirectory"`

	// DataDirectory holds persistent and private data storage within the application's sandbox using internal memory
	// (on Android, if you need to use external memory, use .externalDataDirectory).
	// On iOS, this directory is not synced with iCloud (use .syncedDataDirectory). (iOS, Android, BlackBerry 10, windows)
	DataDirectory string `js:"dataDirectory"`

	// CacheDirectory holds directory for cached data files or any files that your app can re-create easily.
	// The OS may delete these files when the device runs low on storage, nevertheless, apps should not rely on the OS to delete files in here.
	// (iOS, Android, BlackBerry 10, OSX, windows)
	CacheDirectory string `js:"cacheDirectory"`

	// ExternalApplicationStorageDirectory holds application space on external storage. (Android)
	ExternalApplicationStorageDirectory string `js:"externalApplicationStorageDirectory"`

	// ExternalDataDirectory holds where to put app-specific data files on external storage. (Android)
	ExternalDataDirectory string `js:"externalDataDirectory"`

	// ExternalCacheDirectory holds application cache on external storage. (Android)
	ExternalCacheDirectory string `js:"externalCacheDirectory"`

	// ExternalRootDirectory holds external storage (SD card) root. (Android, BlackBerry 10)
	ExternalRootDirectory string `js:"externalRootDirectory"`

	// TempDirectory holds temp directory that the OS can clear at will. Do not rely on the OS to clear this directory;
	// your app should always remove files as applicable. (iOS, OSX, windows)
	TempDirectory string `js:"tempDirectory"`

	// SyncedDataDirectory holds directory holding app-specific files that should be synced (e.g. to iCloud). (iOS, windows)
	SyncedDataDirectory string `js:"syncedDataDirectory"`

	// DocumentsDirectory holds directory holding files private to the app, but that are meaningful to other application
	// (e.g. Office files). Note that for OSX this is the user's ~/Documents directory. (iOS, OSX)
	DocumentsDirectory string `js:"documentsDirectory"`

	// SharedDirectory holds directory holding files globally available to all applications (BlackBerry 10)
	SharedDirectory string `js:"sharedDirectory"`
}

var instance *js.Object

// Dir holds Directories singleton
var Dir *Directories

func init() {
	goco.OnDeviceReady(func() {
		Dir = &Directories{Object: mo()}
	})

}

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("file")
	}
	return instance
}

// Metadata of file
type Metadata struct {
	*js.Object
	ModificationTime time.Time `js:"modificationTime"`
	Size             int       `js:"size"`
}

// FileSystem type
type FileSystem struct {
	*js.Object
	Name string          `js:"name"`
	Root *DirectoryEntry `js:"root"`
}

type FileError struct {
	*js.Object
	Code int `js:"code"`
}

type Flags struct {
	Create    bool
	Exclusive bool
}

// Entry serves as a base type for the FileEntry and DirectoryEntry types,
// which provide features specific to file system entries representing files and directories, respectively.
type Entry struct {
	*js.Object

	// FileSystem object representing the file system in which the entry is located.
	FileSystem *FileSystem `js:"filesystem"`

	// FullPath provides the full, absolute path from the file system's root to the entry;
	// it can also be thought of as a path which is relative to the root directory, prepended with a "/" character.
	FullPath string `js:"fullPath"`

	// IsDirectory returns a boolean which is true if the entry represents a directory; otherwise, it's false.
	IsDirectory bool `js:"isDirectory"`

	// IsFile returns a boolean which is true if the entry represents a file. If it's not a file, this value is false.
	IsFile bool `js:"isFile"`

	// Name returns a string containing the name of the entry (the final part of the path, after the last "/" character).
	Name string `js:"name"`
}

type DirectoryEntry struct {
	*Entry
}

type FileEntry struct {
	*Entry
}

func (fl *Flags) jsObject() *js.Object {
	flags := js.Global.Get("Object").New()
	if fl != nil {
		flags.Set("create", fl.Create)
		flags.Set("exclusive", fl.Exclusive)
	} else {
		flags.Set("create", false)
		flags.Set("exclusive", false)
	}
	return flags
}

// AsDirectoryEntry wraps Entry type into DirectoryEntry type. (Ensure IsDirectory property is true before calling this method)
func (e *Entry) AsDirectoryEntry() *DirectoryEntry {
	return &DirectoryEntry{Entry: e}
}

// AsFileEntry wraps Entry type into FileEntry type. (Ensure IsFile property is true before calling this method)
func (e *Entry) AsFileEntry() *FileEntry {
	return &FileEntry{Entry: e}
}

// GetMetadata obtains metadata about the file, such as its modification date and size.
func (e *Entry) GetMetadata() (res *Metadata, err error) {
	ch := make(chan struct{})
	success := func(md *Metadata) {
		res = md
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	e.Call("getMetadata", success, fail)
	<-ch
	return
}

// GetParent returns a entry representing the entry's parent directory.
func (e *Entry) GetParent() (res *DirectoryEntry, err error) {
	ch := make(chan struct{})
	success := func(de *DirectoryEntry) {
		res = de
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	e.Call("getParent", success, fail)
	<-ch
	return
}

// CopyTo copies the file specified by the entry to a new target location on the file system.
func (e *Entry) CopyTo(target *Entry) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(en *Entry) {
		res = en
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	e.Call("copyTo", target, success, fail)
	<-ch
	return
}

// MoveTo moves the file or directory to a new location on the file system, or renames the file or directory.
func (e *Entry) MoveTo(target *Entry) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(en *Entry) {
		res = en
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	e.Call("moveTo", target, success, fail)
	<-ch
	return
}

// Remove removes the specified file or directory. You can only remove directories which are empty.
func (e *Entry) Remove() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	e.Call("remove", success, fail)
	<-ch
	return
}

// Read returns all directory entries.
func (d *DirectoryEntry) Read() (res []*Entry, err error) {
	ch := make(chan struct{})
	success := func(entries []*Entry) {
		res = entries
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	reader := d.Call("createReader")
	reader.Call("readEntries", success, fail)
	<-ch
	return
}

// GetDirectory returns a DirectoryEntry instance corresponding to a directory contained somewhere within the directory subtree
// rooted at the directory on which it's called.
func (d *DirectoryEntry) GetDirectory(path string, fl *Flags) (res *DirectoryEntry, err error) {
	ch := make(chan struct{})
	success := func(de *DirectoryEntry) {
		res = de
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	d.Call("getDirectory", path, fl.jsObject(), success, fail)
	<-ch
	return
}

// GetFile returns a FileEntry instance corresponding to a file contained somewhere within the directory subtree
// rooted at the directory on which it's called.
func (d *DirectoryEntry) GetFile(path string, fl *Flags) (res *FileEntry, err error) {
	ch := make(chan struct{})
	success := func(fe *FileEntry) {
		res = fe
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	d.Call("getFile", path, fl.jsObject(), success, fail)
	<-ch
	return
}

func (f *FileEntry) createWriter() (res *js.Object, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = ob
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	f.Call("createWriter", success, fail)
	<-ch
	return
}

func (f *FileEntry) file() (res *js.Object) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = ob
		close(ch)
	}
	f.Call("file", success)
	<-ch
	return
}

func (f *FileEntry) Write(data []byte) (err error) {
	writer, err := f.createWriter()
	if err != nil {
		return err
	}
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	writer.Set("onwriteend", success)
	writer.Set("onerror", fail)
	buffer := js.Global.Get("Uint8Array").New(data).Get("buffer")
	blob := js.Global.Get("Blob").New([]*js.Object{buffer}, map[string]interface{}{"type": ""})
	writer.Call("write", blob)
	<-ch
	return
}

func (f *FileEntry) Read() (res []byte, err error) {
	reader := js.Global.Get("FileReader").New()
	blob := f.file()
	ch := make(chan struct{})
	success := func() {
		arrayBuffer := reader.Get("result")
		res = js.Global.Get("Uint8Array").New(arrayBuffer).Interface().([]byte)
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	reader.Set("onloadend", success)
	reader.Set("onerror", fail)
	reader.Call("readAsArrayBuffer", blob)
	<-ch
	return
}

// ToURL creates and returns a URL which identifies the entry.
func (e *Entry) ToURL() (res string) {
	return e.Call("toURL").String()
}

// ResolveLocalFileSystemURL retrieves a Entry instance based on it's local URL
func ResolveLocalFileSystemURL(url string) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Entry{Object: ob}
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	js.Global.Call("resolveLocalFileSystemURL", url, success, fail)
	<-ch
	return
}

// RequestFileSystem requests a file system where data should be stored.
// typ is the storage type of the file system.
// size is the storage space—in bytes—that you need for your app.
func RequestFileSystem(typ int, size int) (res *FileSystem, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &FileSystem{Object: ob}
		close(ch)
	}
	fail := func(e *FileError) {
		err = e.Error()
		close(ch)
	}
	js.Global.Call("requestFileSystem", typ, size, success, fail)
	<-ch
	return
}
