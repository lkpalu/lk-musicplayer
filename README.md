# lk-musicplayer
a small commend music player

# Usages
1. install sqlite3
2. Add the directory containing lk.exe to the system environment variable PATH
3. Add the directory containing lk.exe to the system environment variable Mytool

# Commend
-h --help Help document
add --add music
rm  --rm music
show --show music lists
play  --play music

# Example
You can CD it to your music directory and then execute lk add. 
This way, your music will be added to the playlist, 
and you don't need to CD it again the next time you start it. 
You can start it in any directory or terminal and listen to music

```shell
lk add
```

You can execute lk show to see what songs are on the playlist

```shell
lk show
```

You can execute lk rm to delete songs from the playlist, but it only deletes the playlist and does not actually delete the file for that song

```shell
lk rm [n]
```

Most importantly, you can perform lk play to play your song

```shell
lk play [n]
```

