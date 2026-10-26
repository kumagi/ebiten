// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2012 The glfw3-go Authors
// SPDX-FileCopyrightText: 2006-2019 Camilla Löwy <elmindreda@glfw.org>
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd

package glfw

import (
	"fmt"
)

func GetX11Display() (uintptr, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	return _glfw.platformWindow.display, nil
}

// GetX11Adapter returns the RRCrtc of the monitor.
func (m *Monitor) GetX11Adapter() (_RRCrtc, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	return m.platform.crtc, nil
}

// GetX11Monitor returns the RROutput of the monitor.
func (m *Monitor) GetX11Monitor() (_RROutput, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	return m.platform.output, nil
}

// GetX11Window returns the Window of the window.
func (w *Window) GetX11Window() (uintptr, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	return uintptr(w.platform.handle), nil
}

// GetX11InputContext returns the X input context (XIC) of the window, or 0
// when no input method is available.
func (w *Window) GetX11InputContext() (uintptr, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	return w.platform.ic, nil
}

// GetGLXContext returns the GLXContext of the window.
func (w *Window) GetGLXContext() (uintptr, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	if w.context.source != NativeContextAPI {
		return 0, fmt.Errorf("glfw: window has no GLX context: %w", NoWindowContext)
	}
	return w.context.platform.glx.handle, nil
}

// GetGLXWindow returns the GLXWindow of the window.
func (w *Window) GetGLXWindow() (uintptr, error) {
	if !_glfw.initialized {
		return 0, NotInitialized
	}
	if w.context.source != NativeContextAPI {
		return 0, fmt.Errorf("glfw: window has no GLX context: %w", NoWindowContext)
	}
	return uintptr(w.context.platform.glx.window), nil
}

// SetX11SelectionString sets the X11 primary selection string.
func SetX11SelectionString(str string) error {
	if !_glfw.initialized {
		return NotInitialized
	}

	_glfw.platformWindow.primarySelectionString = str

	xSetSelectionOwner(_glfw.platformWindow.display,
		_glfw.platformWindow.PRIMARY,
		_glfw.platformWindow.helperWindowHandle,
		_CurrentTime)

	if xGetSelectionOwner(_glfw.platformWindow.display, _glfw.platformWindow.PRIMARY) !=
		_glfw.platformWindow.helperWindowHandle {
		return fmt.Errorf("glfw: x11: failed to become owner of primary selection: %w", PlatformError)
	}
	return nil
}

// GetX11SelectionString gets the X11 primary selection string.
func GetX11SelectionString() (string, error) {
	if !_glfw.initialized {
		return "", NotInitialized
	}
	return getSelectionString(_glfw.platformWindow.PRIMARY)
}
