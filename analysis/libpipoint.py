# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import datetime
import collections
import re

import numpy as np

Doc = collections.namedtuple('Doc', 'events')
Event = collections.namedtuple('Event', 'stamp filename name vtype values')


def parse_primitive(value):
    try:
        return float(value)
    except ValueError:
        return value


def parse_values(vtype, values):
    if vtype == 'float64':
        return float(values)
    elif vtype == 'int':
        return int(values)
    elif vtype.startswith('*'):
        match = re.match(r'&.+{(.+)}$', values)
        assert match
        if match:
            fields = match.group(1).split(',')
            out = collections.OrderedDict()
            for field in fields:
                parts = field.strip().split(':')
                if len(parts) == 2:
                    name, value = parts
                    out[name] = parse_primitive(value)
            return out
        return values
    else:
        assert False, values
        return values


def _parse(lines, predicate=None):
    for line in lines:
        match = re.match(
            r'(\d+)/(\d+)/(\d+) (\d+):(\d+):([\d.]+) (\S+): (\S+) (\S+) (.*)',
            line)
        if match:
            (year, month, day, hour, minute, seconds, filename, name, vtype,
             values) = match.groups()

            if predicate and not predicate(name):
                continue

            seconds = float(seconds)
            stamp = datetime.datetime(
                int(year),
                int(month),
                int(day),
                int(hour),
                int(minute), int(seconds), int((seconds - int(seconds)) * 1e6))

            v = parse_values(vtype, values)
            yield Event(stamp, filename, name, vtype, v)


def parse(lines, predicate=None):
    return Doc(list(_parse(lines, predicate)))


def deriv1(v, dt):
    yield 0
    for i in range(1, len(v)):
        n1, n0 = v[i - 1], v[i]
        yield (n0 - n1) / dt


def deriv2(v, t):
    yield 0
    for i in range(1, len(v)):
        n1, n0 = v[i - 1], v[i]
        t1, t0 = t[i - 1], t[i]
        yield (n0 - n1) / (t0 - t1)


def deriv3(v, t):
    yield 0
    yield 0
    for i in range(2, len(v)):
        n2, n1, n0 = v[i - 2], v[i - 1], v[i]
        t2, t1, t0 = t[i - 2], t[i - 1], t[i]
        yield (1 * n2 - 4 * n1 + 3 * n0) / (1 * t2 - 4 * t1 + 3 * t0)


def deriv4(v, t):
    yield 0
    for i in range(1, len(v) - 1):
        n2, n1, n0 = v[i - 1], v[i - 0], v[i + 1]
        t2, t1, t0 = t[i - 1], t[i - 0], t[i + 1]
        yield (-1 * n2 - 0 * n1 + 1 * n0) / (-1 * t2 - 0 * t1 + 1 * t0)

    yield 0


def find(t, tt):
    earlier = [x for x in t if x <= tt]
    return len(earlier) - 1


def lpred(t, s, v):
    t1 = np.arange(min(t), max(t), 0.02)
    s1 = []

    for tx in t1:
        idx = find(t, tx)
        t0, s0, v0 = t[idx], s[idx], v[idx]
        dt = tx - t0
        s1.append(s0 + v0 * dt)

    return t1, s1
