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
	"os"

	"gobot.io/x/gobot/sysfs"
)

type PwmPin struct {
	Chip int
	Pin  int
}

func writeFile(path string, value int) (wrote int, err error) {
	file, err := sysfs.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write([]byte(fmt.Sprintf("%d\n", value)))
}

func (p *PwmPin) chip() string {
	return fmt.Sprintf("/sys/class/pwm/pwmchip%d", p.Chip)
}

func (p *PwmPin) attr(attr string) string {
	return fmt.Sprintf("/sys/class/pwm/pwmchip%d/pwm%d/%s", p.Chip, p.Pin, attr)
}

func (p *PwmPin) SetEnable(val int) (err error) {
	_, err = writeFile(p.attr("enable"), val)
	return
}

func (p *PwmPin) SetPeriod(period int) (err error) {
	_, err = writeFile(p.attr("period"), period)
	return
}

func (p *PwmPin) SetDuty(duty int) (err error) {
	_, err = writeFile(p.attr("duty_cycle"), duty)
	return
}

func (p *PwmPin) Export() (err error) {
	path := p.chip() + "/export"
	_, err = writeFile(path, p.Pin)
	return
}

func (p *PwmPin) UnExport() (err error) {
	path := p.chip() + "/unexport"
	_, err = writeFile(path, p.Pin)
	return
}
