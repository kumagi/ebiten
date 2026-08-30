// Copyright 2026 The Ebitengine Authors
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

package gamepad

/*
#include <jni.h>
#include <stdint.h>

// Basically the following code is equivalent to the following Java code:
//
//     // InputDevice.getVibratorManager requires the API Level 31 or newer.
//     if (Build.VERSION.SDK_INT < 31) {
//       return
//     }
//     final InputDevice device = InputDevice.getDevice(deviceID)
//     if (device == null) {
//       return
//     }
//     final VibratorManager manager = device.getVibratorManager()
//     final int[] ids = manager.getVibratorIds()
//     if (2 <= ids.length) {
//       // The vibrators are ordered from the most prominent one, which is the low-frequency motor, to
//       // the least prominent one.
//       vibrate(manager.getVibrator(ids[0]), strongAmplitude, milliseconds)
//       vibrate(manager.getVibrator(ids[1]), weakAmplitude, milliseconds)
//     } else if (ids.length == 1) {
//       // A device with a single motor cannot vibrate its motors separately, so mix them. The
//       // low-frequency motor is more noticeable than the high-frequency motor.
//       vibrate(manager.getVibrator(ids[0]), (3*strongAmplitude + 2*weakAmplitude + 2) / 5, milliseconds)
//     }
//
// and `vibrate` is:
//
//     void vibrate(Vibrator vibrator, int amplitude, long milliseconds) {
//       if (amplitude <= 0) {
//         // A zero amplitude stops vibrating.
//         vibrator.cancel()
//         return
//       }
//       int a = VibrationEffect.DEFAULT_AMPLITUDE // -1
//       if (vibrator.hasAmplitudeControl()) {
//         a = amplitude
//       }
//       vibrator.vibrate(VibrationEffect.createOneShot(milliseconds, a))
//     }
//
// Note that this requires a manifest setting:
//
//     <uses-permission android:name="android.permission.VIBRATE"/>
//
#cgo noescape vibrateGamepad
#cgo nocallback vibrateGamepad

// clearException clears a pending Java exception, if any. The exception must be cleared, as otherwise
// it is thrown to the caller of the native function, which is the Go runtime and not a Java method.
static void clearException(JNIEnv* env) {
  if ((*env)->ExceptionCheck(env)) {
    (*env)->ExceptionClear(env);
  }
}

// vibrateVibrator vibrates the vibrator with the amplitude, which is in the range from 0 to 255.
// The zero amplitude stops vibrating the vibrator.
static void vibrateVibrator(JNIEnv* env,
                            const jmethodID vibrateMethod,
                            const jmethodID cancelMethod,
                            const jmethodID hasAmplitudeControlMethod,
                            const jclass vibrationEffectClass,
                            const jmethodID createOneShotMethod,
                            const jobject vibrator,
                            const int amplitude,
                            const int64_t milliseconds) {
  if (!vibrator) {
    return;
  }

  if (amplitude <= 0) {
    (*env)->CallVoidMethod(env, vibrator, cancelMethod);
    clearException(env);
    return;
  }

  // VibrationEffect.DEFAULT_AMPLITUDE (-1) must be used for a vibrator whose amplitude cannot be
  // controlled.
  jint amplitudeToUse = -1;
  if ((*env)->CallBooleanMethod(env, vibrator, hasAmplitudeControlMethod)) {
    amplitudeToUse = amplitude;
  }
  clearException(env);

  const jobject effect = (*env)->CallStaticObjectMethod(
      env, vibrationEffectClass, createOneShotMethod, (jlong)milliseconds, amplitudeToUse);
  if (!effect) {
    clearException(env);
    return;
  }

  // vibrate throws a SecurityException if the VIBRATE permission is not granted.
  (*env)->CallVoidMethod(env, vibrator, vibrateMethod, effect);
  clearException(env);

  (*env)->DeleteLocalRef(env, effect);
}

static void vibrateGamepad(uintptr_t jni_env, int32_t deviceID, int64_t milliseconds, int32_t strongAmplitude, int32_t weakAmplitude) {
  JNIEnv* env = (JNIEnv*)jni_env;

  static int apiLevel = 0;
  if (!apiLevel) {
    const jclass android_os_Build_VERSION = (*env)->FindClass(env, "android/os/Build$VERSION");
    if (!android_os_Build_VERSION) {
      clearException(env);
      return;
    }

    apiLevel = (*env)->GetStaticIntField(
        env, android_os_Build_VERSION,
        (*env)->GetStaticFieldID(env, android_os_Build_VERSION, "SDK_INT", "I"));

    (*env)->DeleteLocalRef(env, android_os_Build_VERSION);
  }

  // A gamepad cannot be vibrated before Android 12 (API Level 31), which is when
  // InputDevice.getVibratorManager was added.
  if (apiLevel < 31) {
    return;
  }

  const jclass android_view_InputDevice = (*env)->FindClass(env, "android/view/InputDevice");
  const jclass android_os_VibratorManager = (*env)->FindClass(env, "android/os/VibratorManager");
  const jclass android_os_Vibrator = (*env)->FindClass(env, "android/os/Vibrator");
  const jclass android_os_VibrationEffect = (*env)->FindClass(env, "android/os/VibrationEffect");

  const jmethodID getDeviceMethod = (*env)->GetStaticMethodID(
      env, android_view_InputDevice, "getDevice", "(I)Landroid/view/InputDevice;");
  const jmethodID getVibratorManagerMethod = (*env)->GetMethodID(
      env, android_view_InputDevice, "getVibratorManager", "()Landroid/os/VibratorManager;");
  const jmethodID getVibratorIdsMethod = (*env)->GetMethodID(
      env, android_os_VibratorManager, "getVibratorIds", "()[I");
  const jmethodID getVibratorMethod = (*env)->GetMethodID(
      env, android_os_VibratorManager, "getVibrator", "(I)Landroid/os/Vibrator;");
  const jmethodID vibrateMethod = (*env)->GetMethodID(
      env, android_os_Vibrator, "vibrate", "(Landroid/os/VibrationEffect;)V");
  const jmethodID cancelMethod = (*env)->GetMethodID(env, android_os_Vibrator, "cancel", "()V");
  const jmethodID hasAmplitudeControlMethod = (*env)->GetMethodID(
      env, android_os_Vibrator, "hasAmplitudeControl", "()Z");
  const jmethodID createOneShotMethod = (*env)->GetStaticMethodID(
      env, android_os_VibrationEffect, "createOneShot", "(JI)Landroid/os/VibrationEffect;");

  jobject device = NULL;
  jobject manager = NULL;
  jintArray ids = NULL;
  jint* vibratorIDs = NULL;
  jobject strongVibrator = NULL;
  jobject weakVibrator = NULL;
  jsize idCount = 0;
  int32_t mixedAmplitude = 0;

  if (!android_view_InputDevice || !android_os_VibratorManager || !android_os_Vibrator || !android_os_VibrationEffect ||
      !getDeviceMethod || !getVibratorManagerMethod || !getVibratorIdsMethod || !getVibratorMethod ||
      !vibrateMethod || !cancelMethod || !hasAmplitudeControlMethod || !createOneShotMethod) {
    clearException(env);
    goto end;
  }

  // getDevice returns null if the device does not exist.
  device = (*env)->CallStaticObjectMethod(env, android_view_InputDevice, getDeviceMethod, (jint)deviceID);
  if (!device) {
    clearException(env);
    goto end;
  }

  // getVibratorManager throws an UnsupportedOperationException if the device cannot vibrate.
  manager = (*env)->CallObjectMethod(env, device, getVibratorManagerMethod);
  if (!manager) {
    clearException(env);
    goto end;
  }

  ids = (jintArray)(*env)->CallObjectMethod(env, manager, getVibratorIdsMethod);
  if (!ids) {
    clearException(env);
    goto end;
  }

  idCount = (*env)->GetArrayLength(env, ids);
  vibratorIDs = (jint*)(*env)->GetIntArrayElements(env, ids, NULL);
  if (!vibratorIDs) {
    clearException(env);
    goto end;
  }

  mixedAmplitude = strongAmplitude;
  if (2 <= idCount) {
    // The vibrators are ordered from the most prominent one, which is the low-frequency motor, to the
    // least prominent one.
    strongVibrator = (*env)->CallObjectMethod(env, manager, getVibratorMethod, vibratorIDs[0]);
    weakVibrator = (*env)->CallObjectMethod(env, manager, getVibratorMethod, vibratorIDs[1]);
  } else if (idCount == 1) {
    strongVibrator = (*env)->CallObjectMethod(env, manager, getVibratorMethod, vibratorIDs[0]);
    // A device with a single motor cannot vibrate its motors separately, so mix the amplitudes. The
    // low-frequency motor is more noticeable than the high-frequency motor.
    mixedAmplitude = (3 * strongAmplitude + 2 * weakAmplitude + 2) / 5;
  }
  clearException(env);
  (*env)->ReleaseIntArrayElements(env, ids, vibratorIDs, JNI_ABORT);
  vibratorIDs = NULL;

  vibrateVibrator(env, vibrateMethod, cancelMethod, hasAmplitudeControlMethod,
                  android_os_VibrationEffect, createOneShotMethod,
                  strongVibrator, mixedAmplitude, milliseconds);
  vibrateVibrator(env, vibrateMethod, cancelMethod, hasAmplitudeControlMethod,
                  android_os_VibrationEffect, createOneShotMethod,
                  weakVibrator, weakAmplitude, milliseconds);

end:
  if (vibratorIDs) {
    (*env)->ReleaseIntArrayElements(env, ids, vibratorIDs, JNI_ABORT);
  }
  if (strongVibrator) {
    (*env)->DeleteLocalRef(env, strongVibrator);
  }
  if (weakVibrator) {
    (*env)->DeleteLocalRef(env, weakVibrator);
  }
  if (ids) {
    (*env)->DeleteLocalRef(env, ids);
  }
  if (manager) {
    (*env)->DeleteLocalRef(env, manager);
  }
  if (device) {
    (*env)->DeleteLocalRef(env, device);
  }
  if (android_os_VibrationEffect) {
    (*env)->DeleteLocalRef(env, android_os_VibrationEffect);
  }
  if (android_os_Vibrator) {
    (*env)->DeleteLocalRef(env, android_os_Vibrator);
  }
  if (android_os_VibratorManager) {
    (*env)->DeleteLocalRef(env, android_os_VibratorManager);
  }
  if (android_view_InputDevice) {
    (*env)->DeleteLocalRef(env, android_view_InputDevice);
  }
}

*/
import "C"

import (
	"github.com/ebitengine/gomobile/app"
)

// vibrateAndroidGamepad vibrates the gamepad of the given Android device ID with the given duration in
// milliseconds and with the given amplitudes in the range from 0 to 255, whose strongAmplitude is for
// the low-frequency motor and weakAmplitude is for the high-frequency motor. A zero amplitude stops
// vibrating the corresponding motor.
//
// vibrateAndroidGamepad does nothing if the gamepad cannot be vibrated, or if the Android API Level is
// older than 31.
func vibrateAndroidGamepad(androidDeviceID int, milliseconds int64, strongAmplitude, weakAmplitude int) {
	// Vibrating a gamepad requires the JVM, and calling the JVM can block, so this is done in another
	// goroutine.
	go func() {
		_ = app.RunOnJVM(func(_, jniEnv, _ uintptr) error {
			C.vibrateGamepad(C.uintptr_t(jniEnv),
				C.int32_t(androidDeviceID),
				C.int64_t(milliseconds),
				C.int32_t(strongAmplitude),
				C.int32_t(weakAmplitude))
			return nil
		})
	}()
}
