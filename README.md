## Introduction
  Mural is a 2D graphical frontend for Flame RPG engine written in Go with Pixel library.

  GUI uses [MTK](https://github.com/Isangeles/mtk), simple toolkit built with Pixel library.

  Currently in a early development stage.

  ### Flame games with Mural support:

### Arena

  Description: simple demo game that presents [Flame engine](https://github.com/isangeles/flame) and Mural GUI features.

  Download: [Linux](https://drive.google.com/open?id=1CAUiHdGq8sxrrNWkRwF1QSaNSVWLKDVg), [Windows](https://drive.google.com/open?id=1rR_k_39o-hqTywUZO628ggA3iN7ZBZTJ)

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

## Controls
### HUD:
WSAD - move HUD camera

Right mouse button - target object

Left mouse button - move player/interact with object(loot/dialog/attack)

SPACE - pause game

ESCAPE - open in-game menu

~ - activate chat

B - open inventory

K - open skills menu

L - open journal

V - open crafting menu

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
[Burn](https://github.com/Isangeles/burn) CI handles commands execution.

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

## Contributing
You are welcome to contribute to project development.

If you looking for things to do, then check TODO file.

When you finish, open pull request to merge your changes with main branch.

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
