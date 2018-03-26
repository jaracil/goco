// Package diagnostic is a GopherJS wrapper for cordova diagnostic plugin.
//
// Install plugin:
//  cordova plugin add cordova.plugins.diagnostic
package diagnostic

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

// Permissions defines Android permissions
type Permissions struct {
	*js.Object
	ReadCalendar         string `js:"READ_CALENDAR"`
	WriteCalendar        string `js:"WRITE_CALENDAR"`
	Camera               string `js:"CAMERA"`
	ReadContacts         string `js:"READ_CONTACTS"`
	WriteContacts        string `js:"WRITE_CONTACTS"`
	GetAccounts          string `js:"GET_ACCOUNTS"`
	AccessFineLocation   string `js:"ACCESS_FINE_LOCATION"`
	AccessCoarseLocation string `js:"ACCESS_COARSE_LOCATION"`
	RecordAudio          string `js:"RECORD_AUDIO"`
	ReadPhoneState       string `js:"READ_PHONE_STATE"`
	CallPhone            string `js:"CALL_PHONE"`
	AddVoicemail         string `js:"ADD_VOICEMAIL"`
	UseSip               string `js:"USE_SIP"`
	ProcessOutgoingCalls string `js:"PROCESS_OUTGOING_CALLS"`
	ReadCallLog          string `js:"READ_CALL_LOG"`
	WriteCallLog         string `js:"WRITE_CALL_LOG"`
	SendSMS              string `js:"SEND_SMS"`
	ReceiveSMS           string `js:"RECEIVE_SMS"`
	ReadSMS              string `js:"READ_SMS"`
	ReceiveWapPush       string `js:"RECEIVE_WAP_PUSH"`
	ReceiveMMS           string `js:"RECEIVE_MMS"`
	WriteExternalStorage string `js:"WRITE_EXTERNAL_STORAGE"`
	ReadExternalStorage  string `js:"READ_EXTERNAL_STORAGE"`
	BodySensors          string `js:"BODY_SENSORS"`
}

// PermissionStatus defines possible permission status
type PermissionStatus struct {
	*js.Object
	Granted          string `js:"GRANTED"`
	Denied           string `js:"DENIED"`
	DeniedAlways     string `js:"DENIED_ALWAYS"`
	NotRequested     string `js:"NOT_REQUESTED"`
	Restricted       string `js:"RESTRICTED"`
	GrantedWhenInUse string `js:"GRANTED_WHEN_IN_USE"`
}

// Architectures defines possible hardware architectures
type Architectures struct {
	*js.Object
	Unknown string `js:"UNKNOWN"`
	ArmV6   string `js:"ARMv6"`
	ArmV7   string `js:"ARMv7"`
	ArmV8   string `js:"ARMv8"`
	X86     string `js:"X86"`
	X86_64  string `js:"X86_64"`
	Mips    string `js:"MIPS"`
	Mips_64 string `js:"MIPS_64"`
}

var (
	instance *js.Object

	// PermStatus is an instance of PermissionStatus
	PermStatus *PermissionStatus
	// Perm is an instance of Permissions
	Perm *Permissions
	// Arch is an instance of Architectures
	Arch *Architectures
)

func init() {
	goco.OnDeviceReady(func() {
		Perm = &Permissions{Object: mo().Get("permission")}
		PermStatus = &PermissionStatus{Object: mo().Get("permissionStatus")}
		Arch = &Architectures{Object: mo().Get("cpuArchitecture")}
	})
}

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("plugins").Get("diagnostic")
	}
	return instance
}

// SwitchToSettings opens settings page for this app.
//  Platforms: Android and iOS
//  Notes:
//   Android: this opens the "App Info" page in the Settings app.
//   iOS: this opens the app settings page in the Settings app. This works only on iOS 8+ - iOS 7 and below will return with error.
func SwitchToSettings() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("switchToSettings", success, fail)
	<-ch
	return
}

// SwitchToWirelessSettings Switches to the wireless settings page in the Settings app. Allows configuration of wireless controls such as Wi-Fi, Bluetooth and Mobile networks.
//  Platforms: Android
func SwitchToWirelessSettings() (err error) {
	mo().Call("switchToWirelessSettings")
	err = nil
	return
}

// SwitchToMobileDataSettings Displays mobile settings to allow user to enable mobile data.
//  Platforms: Android and Windows 10 UWP
func SwitchToMobileDataSettings() (err error) {
	mo().Call("switchToMobileDataSettings")
	err = nil
	return
}

// GetPermissionAuthorizationStatus returns the status of permission
//  Platforms: Android
func GetPermissionAuthorizationStatus(perm string) (stat string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		stat = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getPermissionAuthorizationStatus", success, fail, perm)
	<-ch
	return
}

// GetPermissionsAuthorizationStatus returns the status of permissions
//  Platforms: Android
func GetPermissionsAuthorizationStatus(perms []string) (stat map[string]string, err error) {
	ch := make(chan struct{})
	success := func(st map[string]interface{}) {
		stat = make(map[string]string)
		for k, v := range st {
			stat[k] = v.(string)
		}
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getPermissionsAuthorizationStatus", success, fail, perms)
	<-ch
	return
}

// RequestRuntimePermission requests app to be granted authorization for a runtime permission.
//  Platforms: Android
//  Note: this is intended for Android 6 / API 23 and above. Calling on Android 5 / API 22 and below
//  will have no effect as the permissions are already granted at installation time.
func RequestRuntimePermission(perm string) (stat string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		stat = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("requestRuntimePermission", success, fail, perm)
	<-ch
	return
}

// RequestRuntimePermissions requests app to be granted authorization for multiple runtime permissions.
//  Platforms: Android
//  Note: this is intended for Android 6 / API 23 and above. Calling on Android 5 / API 22 and below
//  will always return GRANTED status as permissions are already granted at installation time.
func RequestRuntimePermissions(perms []string) (stat map[string]string, err error) {
	ch := make(chan struct{})
	success := func(st map[string]interface{}) {
		stat = make(map[string]string)
		for k, v := range st {
			stat[k] = v.(string)
		}
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("requestRuntimePermissions", success, fail, perms)
	<-ch
	return
}

// IsRequestingPermission indicates if the plugin is currently requesting a runtime permission via the native API.
//  Platforms: Android
func IsRequestingPermission() bool {
	return mo().Call("isRequestingPermission").Bool()
}

// IsDataRoamingEnabled checks if the device data roaming setting is enabled. Returns true if data roaming is enabled.
//  Platforms: Android
func IsDataRoamingEnabled() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isDataRoamingEnabled", success, fail)
	<-ch
	return
}

// IsADBModeEnabled checks if the device setting for ADB(debug) is switched on. Returns true if ADB(debug) setting is switched on.
//  Platforms: Android
func IsADBModeEnabled() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isADBModeEnabled", success, fail)
	<-ch
	return
}

// IsDeviceRooted checks if the device is rooted. Returns true if the device is rooted.
//  Platforms: Android
func IsDeviceRooted() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isDeviceRooted", success, fail)
	<-ch
	return
}

// IsBackgroundRefreshAuthorized checks if the application is authorized for background refresh.
//  Platforms: iOS
func IsBackgroundRefreshAuthorized() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isBackgroundRefreshAuthorized", success, fail)
	<-ch
	return
}

// GetBackgroundRefreshStatus returns the background refresh authorization status for the application.
//  Platforms: iOS
func GetBackgroundRefreshStatus() (res string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getBackgroundRefreshStatus", success, fail)
	<-ch
	return
}

// GetArchitecture returns the CPU architecture of the current device.
//  Platforms: Android and iOS
func GetArchitecture() (res string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getArchitecture", success, fail)
	<-ch
	return
}

// Restart restarts the application. By default, a "warm" restart will be performed in which the main Cordova activity
// is immediately restarted, causing the Webview instance to be recreated.
//
// However, if the cold parameter is set to true, then the application will be "cold" restarted, meaning a system exit
// will be performed, causing the entire application to be restarted. This is useful if you want to fully reset the native
// application state but will cause the application to briefly disappear and re-appear.
//  Platforms: Android
func Restart(cold bool) (err error) {
	ch := make(chan struct{})
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("restart", fail, cold)
	<-ch
	return
}

// EnableDebug enables debug mode, which logs native debug messages to the native and JS consoles.
//  Platforms: Android and iOS
func EnableDebug() {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("enableDebug", success)
	<-ch
	return
}
