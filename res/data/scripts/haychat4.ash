# Script that sends arg3 text on chat channel of character with
# arg 2 serial ID if raw distance between him and character with
# arg 1 serial ID is less than 50, after that waits 5 secs.
@1 = testchar#0
@2 = "hay you!"
true {
	for(@3 = moduleshow -o area-chars -t area1_main) {
     		rawdis(@1, @3) < 50 {
        		charman -o set -a chat @2 -t @1;
        		wait(5);
	     	};
	};
}
