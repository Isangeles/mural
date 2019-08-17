# Script that sends arg2 text on chat channel of character with
# arg1 serial ID if raw distance between him and any character from
# area with ID 'area1_main'(expect char with arg1 serial ID) is less 
# than 50, after that script halts for 5 secs.
@1 = testchar#0
@2 = out(chaptershow -o lang -a cdTestcharHay1)
true {
	for(@3 = out(moduleshow -o area-chars -t area1_main)) {
		@1 != @3 {
	     		rawdis(@1, @3) < 50 {
				charman -o set -a chat '@2' -t @1;
				wait(5);
		     	};
		};
	};
};
