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
import argparse
import fileinput
import pickle

import libpipoint


def main():
    parser = argparse.ArgumentParser(description='Parse elog files.')
    parser.add_argument('input', nargs='+')
    parser.add_argument('-o', help='Output filename', required=True)

    drop = set(('pantilt.pan.pv', 'pantilt.tilt.pv', 'pred.position',
                'tick', ))

    args = parser.parse_args()

    for name in args.input:
        lines = fileinput.input(name)
        doc = libpipoint.parse(lines, lambda x: x not in drop)

        with open(args.o, 'wb') as f:
            pickle.dump(doc, f)


if __name__ == '__main__':
    main()
