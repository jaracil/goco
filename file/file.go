// Package file is a GopherJS wrapper for cordova file plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-file
package file

import (
	"errors"
	"time"

	"github.com/gopherjs/gopherjs/js"
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
	ErrUnknown               = errors.New("Unknown error")
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

func code2Err(code int) error {
	switch code {
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
	return ErrUnknown
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("file")
	}
	return instance
}

type Metadata struct {
	*js.Object
	ModificationTime time.Time `js:"modificationTime"`
	Size             int       `js:"size"`
}

type FileSystem struct {
	*js.Object
	Name string `js:"name"`
	Root *Entry `js:"root"`
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

// GetMetadata obtains metadata about the file, such as its modification date and size.
func (e *Entry) GetMetadata() (res *Metadata, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Metadata{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	e.Call("getMetadata", success, fail)
	<-ch
	return
}

// GetParent returns a entry representing the entry's parent directory.
func (e *Entry) GetParent() (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Entry{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	e.Call("getParent", success, fail)
	<-ch
	return
}

// CopyTo copies the file specified by the entry to a new target location on the file system.
func (e *Entry) CopyTo(target *Entry) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Entry{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	e.Call("copyTo", target, success, fail)
	<-ch
	return
}

// MoveTo moves the file or directory to a new location on the file system, or renames the file or directory.
func (e *Entry) MoveTo(target *Entry) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Entry{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
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
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	e.Call("remove", success, fail)
	<-ch
	return
}

// ToURL creates and returns a URL which identifies the entry.
func (e *Entry) ToURL() (res string) {
	return e.Call("toURL").String()
}

// ApplicationDirectory returns read-only directory where the application is installed. (iOS, Android, BlackBerry 10, OSX, windows)
func ApplicationDirectory() string {
	return mo().Get("applicationDirectory").String()
}

// ApplicationStorageDirectory returns root directory of the application's sandbox;
// on iOS & windows this location is read-only (but specific subdirectories [like /Documents on iOS or /localState on windows] are read-write).
// All data contained within is private to the app. (iOS, Android, BlackBerry 10, OSX)
func ApplicationStorageDirectory() string {
	return mo().Get("applicationStorageDirectory").String()
}

// DataDirectory returns persistent and private data storage within the application's sandbox using internal memory
// (on Android, if you need to use external memory, use .externalDataDirectory).
// On iOS, this directory is not synced with iCloud (use .syncedDataDirectory). (iOS, Android, BlackBerry 10, windows)
func DataDirectory() string {
	return mo().Get("dataDirectory").String()
}

// CacheDirectory returns directory for cached data files or any files that your app can re-create easily.
// The OS may delete these files when the device runs low on storage, nevertheless, apps should not rely on the OS to delete files in here.
// (iOS, Android, BlackBerry 10, OSX, windows)
func CacheDirectory() string {
	return mo().Get("cacheDirectory").String()
}

// ExternalApplicationStorageDirectory returns application space on external storage. (Android)
func ExternalApplicationStorageDirectory() string {
	return mo().Get("externalApplicationStorageDirectory").String()
}

// ExternalDataDirectory returns where to put app-specific data files on external storage. (Android)
func ExternalDataDirectory() string {
	return mo().Get("externalDataDirectory").String()
}

// ExternalCacheDirectory returns application cache on external storage. (Android)
func ExternalCacheDirectory() string {
	return mo().Get("externalCacheDirectory").String()
}

// ExternalRootDirectory returns external storage (SD card) root. (Android, BlackBerry 10)
func ExternalRootDirectory() string {
	return mo().Get("externalRootDirectory").String()
}

// TempDirectory returns temp directory that the OS can clear at will. Do not rely on the OS to clear this directory;
// your app should always remove files as applicable. (iOS, OSX, windows)
func TempDirectory() string {
	return mo().Get("tempDirectory").String()
}

// SyncedDataDirectory returns directory holding app-specific files that should be synced (e.g. to iCloud). (iOS, windows)
func SyncedDataDirectory() string {
	return mo().Get("syncedDataDirectory").String()
}

// DocumentsDirectory returns directory holding files private to the app, but that are meaningful to other application
// (e.g. Office files). Note that for OSX this is the user's ~/Documents directory. (iOS, OSX)
func DocumentsDirectory() string {
	return mo().Get("documentsDirectory").String()
}

// SharedDirectory returns directory holding files globally available to all applications (BlackBerry 10)
func SharedDirectory() string {
	return mo().Get("sharedDirectory").String()
}

func ResolveLocalFileSystemURL(url string) (res *Entry, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &Entry{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	js.Global.Call("resolveLocalFileSystemURL", url, success, fail)
	<-ch
	return
}

// RequestFileSystem requests a file system where data should be stored.
// typ is the storage type of the file system.
// size is the storage space—in bytes—that you need for your app.
func RequestFileSystem(typ string, size int) (res *FileSystem, err error) {
	ch := make(chan struct{})
	success := func(ob *js.Object) {
		res = &FileSystem{Object: ob}
		close(ch)
	}
	fail := func(code int) {
		err = code2Err(code)
		close(ch)
	}
	js.Global.Call("requestFileSystem", typ, size, success, fail)
	<-ch
	return
}
