// Copyright 2020 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphicscommand

import (
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2/internal/thread"
)

// Publish the render thread independently of the goroutine running its loop.
var theRenderThread atomic.Pointer[thread.Thread]

// theNoopRenderThread is the render thread before an OS thread is set as the render thread.
var theNoopRenderThread thread.Thread = thread.NewNoopThread()

func renderThread() thread.Thread {
	if p := theRenderThread.Load(); p != nil {
		return *p
	}
	return theNoopRenderThread
}

// SetOSThreadAsRenderThread sets an OS thread as the rendering thread e.g. for OpenGL.
func SetOSThreadAsRenderThread() {
	t := thread.Thread(thread.NewOSThread())
	theRenderThread.Store(&t)
}

// LoopRenderThread runs the rendering thread loop until done receives a value or is closed.
func LoopRenderThread(done <-chan struct{}) {
	renderThread().Loop(done)
}

// runOnRenderThread calls f with arg on the rendering thread and returns its result.
func runOnRenderThread[A, R any](f func(A) R, arg A) R {
	return thread.CallWithArgAndResult(renderThread(), f, arg)
}

// runOnRenderThreadAsync queues f with arg on the rendering thread.
func runOnRenderThreadAsync[A any](f func(A), arg A) {
	// As the render thread's task channel is unbuffered,
	// CallAsync blocks when the previously-queued task is not executed yet.
	// This blocking is expected as double-buffering is used.
	thread.CallAsync(renderThread(), f, arg)
}

// AbortReadPixels completes pending and queued pixel read-backs with an error.
// Their results have been delivered when AbortReadPixels returns.
// The caller must hold the atlas backend mutex.
func AbortReadPixels() {
	theCommandQueueManager.abortReadPixels(errReadPixelsAborted)
}
