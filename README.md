# Inspiration
Vastly inspired on great YouTube channel [JustForFunc](https://www.youtube.com/channel/UC_BzFbxG2za3bp5NRRRXJSw) !!

# Requirements
Follow system requirements at [https://github.com/veandco/go-sdl2](https://github.com/veandco/go-sdl2)

## OSX
For OSX you are good to go with these brew packages:
```
$ brew install sdl2{,_image,_ttf,_mixer}
```

# Install dependencies
```
$ export GOPATH=`pwd`/vendor:`pwd`
$ make deps
```

# Build
```
$ make build
```

# Run
```
$ bin/flappy
```