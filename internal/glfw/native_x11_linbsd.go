// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2012 The glfw3-go Authors
// SPDX-FileCopyrightText: 2006-2019 Camilla Löwy <elmindreda@glfw.org>
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd

package glfw

import (
	"fmt"
	"unsafe"
)

func GetX11Display() (uintptr, error) {
	return _glfw.platformWindow.display, nil
}

// GetX11Adapter returns the RRCrtc of the monitor.
func (m *Monitor) GetX11Adapter() (_RRCrtc, error) {
	return m.platform.crtc, nil
}

// GetX11CrtcSize returns the current pixel size of crtc, or false if it is unavailable.
func GetX11CrtcSize(crtc _RRCrtc) (width, height int, ok bool) {
	randr := &_glfw.platformWindow.randr
	if !randr.available || _glfw.platformWindow.display == 0 || crtc == 0 {
		return 0, 0, false
	}
	resources := randr.GetScreenResourcesCurrent(_glfw.platformWindow.display, _glfw.platformWindow.root)
	if resources == 0 {
		return 0, 0, false
	}
	defer randr.FreeScreenResources(resources)
	info := getCrtcInfoX11(resources, crtc)
	if info == 0 {
		return 0, 0, false
	}
	defer randr.FreeCrtcInfo(info)
	ci := (*_XRRCrtcInfo)(unsafe.Pointer(info))
	return int(ci.Width), int(ci.Height), true
}

// GetX11Monitor returns the RROutput of the monitor.
func (m *Monitor) GetX11Monitor() (_RROutput, error) {
	return m.platform.output, nil
}

// GetX11Window returns the Window of the window.
func (w *Window) GetX11Window() (uintptr, error) {
	return uintptr(w.platform.handle), nil
}

// GetX11InputContext returns the X input context (XIC) of the window, or 0
// when no input method is available.
func (w *Window) GetX11InputContext() (uintptr, error) {
	return w.platform.ic, nil
}

// GetGLXContext returns the GLXContext of the window.
func (w *Window) GetGLXContext() (uintptr, error) {
	if w.context.source != NativeContextAPI {
		return 0, fmt.Errorf("glfw: window has no GLX context: %w", NoWindowContext)
	}
	return w.context.platform.glx.handle, nil
}

// GetGLXWindow returns the GLXWindow of the window.
func (w *Window) GetGLXWindow() (uintptr, error) {
	if w.context.source != NativeContextAPI {
		return 0, fmt.Errorf("glfw: window has no GLX context: %w", NoWindowContext)
	}
	return uintptr(w.context.platform.glx.window), nil
}

// SetX11SelectionString sets the X11 primary selection string.
func SetX11SelectionString(str string) {
	_glfw.platformWindow.primarySelectionString = str

	xSetSelectionOwner(_glfw.platformWindow.display,
		_glfw.platformWindow.PRIMARY,
		_glfw.platformWindow.helperWindowHandle,
		_CurrentTime)
}

// GetX11SelectionString gets the X11 primary selection string.
func GetX11SelectionString() string {
	str, _ := getSelectionString(_glfw.platformWindow.PRIMARY)
	return str
}
