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
import sys
import pickle
import matplotlib.pylab as plt
import numpy as np

with open(sys.argv[1], 'rb') as f:
    doc = pickle.load(f)

gps = [x for x in doc.events if x.name == 'gps.position']


def flatten(v):
    dim = (len(v), 1 + len(v[0].values))
    out = np.zeros(dim, dtype=np.float64)

    out[:, 0] = [x.stamp.timestamp() for x in v]

    for i, event in enumerate(v):
        out[i, 1:] = list(event.values.values())

    return out


G = flatten(gps)
t = G[:, 1]
dt = np.gradient(t)
lat = G[:, 2] * 110e3
dlat = np.gradient(lat)
v = lat[3:] - lat[:-3]
alt = G[:, 4]

plt.scatter(t, alt - min(alt))
plt.twinx()
plt.scatter(t, lat - min(lat), c='r')
# plt.scatter(t, dt)

# plt.hist(dt, bins=np.arange(0.05, 1.5, 0.1), normed=True)
plt.show()
