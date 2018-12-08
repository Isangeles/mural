## Introduction

  Mural is 2D graphical frontend for Flame RPG engine written in Go with Pixel library.

  Currently in a very early development stage.
  
## Install & Run

  Get [Pixel](https://github.com/faiface/pixel), [go-tmx](https://github.com/salviati/go-tmx/tree/master/tmx) and [Flame](https://github.com/Isangeles/flame).

  Get sources from git:
```
$ go get github.com/isangeles/mural
```

  Build GUI:
```
$ go build github.com/isangeles/mural
```

Copy 'data' directory to directory with Mural executable.

Now, specify the path to a valid Flame module in Flame configuration file,
create file '.flame' in Mural executable directory(or run Mural to create it
automatically) and add the following line:
```
  module:[module name];[module path](optional);
```
If no path provided, the engine will search default modules directory(data/modules).


  Run Mural:
```
$ ./mural
```

## Contact
* Isangeles <<dev@isangeles.pl>>

## License

Copyright 2018 Dariusz Sikora <<dev@isangeles.pl>>
 
This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.
 
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
 
You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
MA 02110-1301, USA.