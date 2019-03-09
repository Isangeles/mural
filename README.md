## Introduction

  Mural is a 2D graphical frontend for Flame RPG engine written in Go with Pixel library.

  GUI uses [MTK](https://github.com/Isangeles/mural/tree/master/core/mtk), simple toolkit built with Pixel library.

  Currently in a very early development stage.
  
## Build

  Get [Pixel](https://github.com/faiface/pixel), [Beep](https://github.com/faiface/beep), [go-tmx](https://github.com/salviati/go-tmx/tree/master/tmx) and [Flame](https://github.com/Isangeles/flame).

  Get sources from git:
```
$ go get -u github.com/isangeles/mural
```

  Build GUI:
```
$ go build github.com/isangeles/mural
```

Copy 'data' directory from /res to directory with Mural executable.

Now, specify the path to a valid Flame module in Flame configuration file.

Create file '.flame' in Mural executable directory(or run Mural to create it
automatically) and add the following line:
```
  module:[module name];[module path](optional);
```
If no path provided, the engine will search default modules directory(data/modules).

Flame modules are available for download [here](http://flame.isangeles.pl/mods).

  Run Mural:
```
$ ./mural
```

## Configuration
Configuration values are loaded from '.mural' file in Mural executable directory.

### Configuration values:
```
  fullscreen:[true/false];
```
Description: enables fullscreen mode, 'true' enables fullscreen, everything else sets windowed mode.

```
  resolution:[width]x[height];
```
Description: specifies current resolution.

```
  map_fow:[true/false];
```
Description: enables 'Fog of War' effect for area map, 'true' enables FOW, everything else sets FOW disabled.

```
  main_font:[file name];
```
Description: specifies name of font file(located in graphic archive) for main UI font.

```
  menu_music:[file name];
```
Description: specifies name of audio file(located in audio archive) for main menu music theme.

```
  button_click_sound:[file name];
```
Description: specifies name of audio file(located in audio archive) for button click sound.

## Commands
[Burn](https://github.com/Isangeles/flame/tree/master/cmd/burn) CI handles commands execution.

Additionally to Burn tools, Mural implements guiman tool to manage game GUI.
You can access CI by the dropdown console in the main menu or chat window in HUD,
both accessible by pressing '`'(grave).

Note: all commands entered in HUD chat window must be prefixed by '$' character.

  Exit mural:
```
guiman -o exit
```
Description: exits program.

  Save HUD state:
```
guiman -o save -t gui-state -a [save name]
```
Description: saves current HUD state to file in current /savegames directory(/savegames/[module]).

  Load HUD state:
```
guiman -o load -t gui-state -a [save name]
```
Description: load HUD state from file in current /savegames directory(/savegames/[module]).


## Contact
* Isangeles <<dev@isangeles.pl>>

## License

Copyright 2018-2019 Dariusz Sikora <<dev@isangeles.pl>>
 
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