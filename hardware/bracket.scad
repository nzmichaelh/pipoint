// The servo mount bracket.
module bracket() {
       t = 3;
       w = 30;
       h = 52;
       d = 36;
       sd = 25;

       $fn = 40;

       // Main arm.
       module arm() {
	 hull () {
	   cube([t, w, w]);
	   translate([0, d, w/2])
	     rotate(90, [0, 1, 0])
	     cylinder(t, w/2, w/2);
	 }
       }

       // Holes in the left arm.
       module servo_holes() {
	   // Servo hole.
	   translate([0, d, w/2])
	     rotate(90, [0,1,0])
	     {
	       // Servo horn indent.
	       translate([0, 0, -1])
		 cylinder(1+1, sd/2, sd/2);
	       // Screw hole.
	       translate([0, 0, -1])
		 cylinder(t+2, 8/2, 8/2);
	       // Outer screw holes.
	       for (i = [0:90:360]) {
		 rotate(i)
		 translate([0, 8, -1])
		   cylinder(t+2, 2.5/2, 2.5/2);
	     }
	   }
       }

       module left_arm() {
	 difference () {
	   arm();
	   servo_holes();
	 }
       }

       module right_arm() {
	 difference () {
	   union () {
	     arm();

	     translate([2, d, w/2])
	     rotate(90, [0,1,0])
	       cylinder(2, sd/3, sd/3);
	   }
           // Bolt hole.
	   translate([2-3, d, w/2])
	     rotate(90, [0,1,0])
	     cylinder(t+2+2, 3/2, 3/2);
	 }
       }

       // Bracket made of two arms and a joiner.
       union () {
	 translate([h, 0, 0])
	   left_arm();

	 translate([0, 0, 0])
	   right_arm();

         cube([h, t, w]);

	 // Fill in the corners.
	 translate([0+t/2, t, 0])
             cylinder(w, t, t);
	 translate([h+t/2, t, 0])
             cylinder(w, t, t);
       }
}
