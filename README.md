# Gopher-Garden
## Overview:
[![License](https://img.shields.io/badge/license-BSD--3--Clause-red)](https://github.com/enzo-sa/gopher-garden/blob/master/LICENSE)
### Play as Garry the Gopher and collect as many carrots as you can whilst avoiding the evil overlords of the garden: _**the snakes**_!
![gopher-garden](https://github.com/enzo-sa/gopher-garden/blob/master/gopher-garden.png)
* Gopher-Garden is a 2D game in which you control Garry the Gopher who moves around a garden and who's goal is to eat as many carrots as he can whilst avoiding deadly snakes which are in constant roam around the garden.

* Gopher-Garden is written entirely in Go 1.14 and uses [Gio](https://gioui.org/), a portable immediate mode GUI implemented in Go and developed by [Elias Naur](https://eliasnaur.com/), for graphics and user-input.

* Gopher-Garden will save all highscores and the user-names that correspond to them in a local file who's data can be viewed through the in-game menu. Do not edit the highscores.txt file because it can erroneate the parsing and cause the program to crash.

## Game Controls:
   ### Text Fields:
*	**At the start of every new round, the user will be prompted to enter their name.**
	* <sub>This name will determine the name assigned to the score they achieve in that game. If their name is not already stored with a highscore in the local file, that name and score pair will be stored in the file and any later scores higher than that score with the same name will update the highscore for that user.</sub>
   ### Key Presses:
**Note:** w-a-s-d or the arrow keys may be used to move Garry.
*	 _**'↑'**_ <sub>or</sub> _**'w'**_   -> Move Garry up.
*	 _**'←'**_ <sub>or</sub> _**'a'**_  -> Move Garry left.
*	 _**'↓'**_ <sub>or</sub> _**'s'**_  -> Move Garry down.
*	 _**'→'**_ <sub>or</sub> _**'d'**_  -> Move Garry right.
*	 _**'Space'**_ <sub>or</sub> _**'e'**_ -> Move Garry into a gopher hole.
	 
	 <sub>When Garry moves into a gopher hole he will pop out at a random one of the other gopher holes. Garry can only move into a gopher hole if his last move was not a gopher hole move.</sub>
   ### Buttons:
* 	_**'Menu'**_  -> On click this button will take you from your current game round into the menu and pause your current game round. <sub> **Note**: If you access the menu in the middle of your game round and you do not resume to that game, your score for that game will NOT be recorded. The menu button can be accessed at all times. When you die however, the continue button acts as an alternative to the menu button which saves your game data.</sub>

	<sub>(The following buttons can only be accessed while in the menu.)</sub>
	*	_**'New Game'**_   -> On click this button will start a new round of the game. 
	*	_**'Back'**_   -> On click this button will take you back into the game round from the menu. 
	*	_**'Highscores'**_ -> On click this button will show the top 5 locally saved highscores and the players who achieved them.
		* 	<sub>The highscores view screen can be exited back into the menu by clicking the _**Back**_ button in the top right of the highscores tab. Note that this back button is different from the menu back button.</sub>
	*	_**'Exit'**_  -> On click this button will terminate the current game session. 
	
	
*	_**'Continue'**_ -> On click this button will take you to the menu and save your score of the game with the name you entered at the beginning of the round. This button is only accessible once you die as an alternative to menu, however pressing menu and not pressing continue will not save your data to the local highscores file.

### Installation:
For a complete installation of Gopher-Garden, please perform the following steps.

*	If the latest version of Go is not already installed on your machine, please follow the system-specific instructions to do so [here](https://golang.org/doc/install).

*	Once the latest version of Go is installed, please install the system-specific dependencies to run [Gio](https://gioui.org/) applications [here](https://gioui.org/#installation).

*	Once these preliminary steps are complete, Gopher-Garden can be run by initializing a go module file in a new folder and running the code hosted on github.	
	*	The module initialization can be done as so: <pre><code>go mod init gopher-garden</pre></code> **Note:** This command must be run in the directory in which you want to initialize the module.
	*	Next, you may run the main Go package from github while in the same directory you initialized the module file in: <pre><code>go run github.com/enzo-sa/gopher-garden/main</pre></code>
Following these steps should allow you to run the game, however if you come across any issues, feel free to notify me through github.

## CREDITS:
This game is based off of the [Go gopher](https://blog.golang.org/gopher#:~:text=The%20Go%20gopher%20is%20an,radio%20station%20in%20New%20Jersey.) which was designed by [Renee French](http://reneefrench.blogspot.com/) and is licensed under [Creative Commons Attribution 3.0](https://creativecommons.org/licenses/by/3.0/legalcode).

A detailed list of image credits can be found in [CREDITS](https://github.com/enzo-sa/gopher-garden/blob/master/ui/resources/CREDITS.md).
## LICENSE:
Gopher-Garden is licensed under __**BSD-3-Clause License**__. For more information see [LICENSE](https://github.com/enzo-sa/gopher-garden/blob/master/LICENSE).
