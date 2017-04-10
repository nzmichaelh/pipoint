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
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"time"
)

// EventLogger is a async event logger.
type EventLogger struct {
	logger *log.Logger
	sink   *os.File
	zip    *gzip.Writer
	msgs   chan []byte
}

// NewEventLogger creates a new event logger that writes to the given
// base name.
func NewEventLogger(name string) *EventLogger {
	now := time.Now().Format(time.RFC3339)
	fname := fmt.Sprintf("%s-%s.txt.gz", name, now)
	sink, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Panicln(err)
	}

	zip := gzip.NewWriter(sink)

	msgs := make(chan []byte, 10)

	el := &EventLogger{
		sink: sink,
		zip:  zip,
		msgs: msgs,
	}

	el.logger = log.New(el, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	go el.run()

	return el
}

func (el *EventLogger) run() {
	tick := time.Tick(3 * time.Second)

	for {
		select {
		case p := <-el.msgs:
			el.zip.Write(p)
		case <-tick:
			el.zip.Flush()
			break
		}
	}
}

func (el *EventLogger) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	el.msgs <- buf
	return len(p), nil
}
