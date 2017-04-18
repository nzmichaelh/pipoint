// A Turnigy 2K actioncam.
module turnigy2k () {
    body=[59, 41, 25];

    //Body.
    cube(body);

    ld = 5;
    lr = 21/2;
    lo = 4.5;
  
    // Lens.
    translate(body + [-lr-lo, -lr-lo, 0])
        cylinder(ld, lr, lr);

    bd = 0.5;
    br = 10/2;

    // OK button.
    translate([6+br, body[1], 7+br])
        rotate(-90, [1, 0, 0])
        cylinder(bd, br, br);

    // Mode button.
    translate([6+br, body[1]-br-10, body[2]])
        cylinder(bd, br, br);

    // Up/down buttons.
    updown = [5, 23, 1];

    //
    translate([0, (body-updown)[1]/2, 6])
        rotate(90, [0, -1, 0])
        cube(updown);

    screen=[31, 23, 10];

    // Screen.
    translate([(body-screen)[0]/2, 12, 0.1-10])
        cube(screen);
}
