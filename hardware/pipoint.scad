include <bracket.scad>;
include <turnigy2k.scad>;

$fn=20;

// Plate that the camera mounts into.
module plate () {
  m=6;
  d=7;
  // Corner radius.
  r=2;
  w=59-r*2+m;
  h=42-r*2+m;

  translate([r, r, 0])
    hull () {
      cylinder(d, r, r);
      translate([w, h, 0])
          cylinder(7, r, r);
      translate([0, h, 0])
          cylinder(7, r, r);
      translate([w, 0, 0])
          cylinder(7, r, r);
  }
  // Text on the top.
  color("grey")
      translate([r*2, h+r*2, d*3/4])
      rotate(90, [-1, 0, 0])
      linear_extrude(0.5)
      text("pipoint v1", size=d/2+1);
}

// Camera mount assembly.
module camera_mount() {
    difference () {
        union () {
            bracket();
            translate([-6, 2, 0])
                rotate(90, [1, 0, 0])
                plate();
        }
        translate([-3, -1, 3])
            rotate(90, [1, 0, 0])
            turnigy2k();
    }
}

camera_mount();
