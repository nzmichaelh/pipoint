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
