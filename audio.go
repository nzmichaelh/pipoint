// Copyright 2017 Google Inc.
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
//

package pipoint

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// AudioOut can play files or speech.
type AudioOut struct {
	queued chan *exec.Cmd
}

// NewAudioOut creates a new, running audio output.
func NewAudioOut() *AudioOut {
	a := &AudioOut{
		queued: make(chan *exec.Cmd, 10),
	}
	go a.run()
	return a
}

// Play plays an audio file.
func (a *AudioOut) Play(path string) {
	a.queued <- exec.Command("ogg123", path)
}

// Say plays a pre-recorded phrase, or falls back to espeak.
func (a *AudioOut) Say(text string) {
	rendered := fmt.Sprintf("phrase/%s.ogg", NormText(text))
	fi, err := os.Stat(rendered)

	if err == nil && fi.Mode().IsRegular() {
		a.Play(rendered)
	} else {
		a.queued <- exec.Command("espeak", text)
	}
}

// run executes the commands to play the sounds.
func (a *AudioOut) run() {
	done := make(chan bool)

	for {
		cmd := <-a.queued
		log.Println(cmd)
		go func() {
			cmd.Run()
			done <- true
		}()

		for c := false; !c; {
			select {
			case <-done:
				log.Println(cmd, "done")
				c = true
			case <-a.queued:
				// Discard anything that comes in
				// while this command is running.
				break
			}
		}
	}
}
