# README

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

- For MacOS:
```
GOOS=darwin GOARCH=amd64 wails build -tags webkit2_41
```

# go-translators with qt5.13.0

## Installation

### Prerequisites
- Essential
```
sudo apt update
sudo apt install -y build-essential libgl1-mesa-dev libx11-dev libxext-dev libxrender-dev libxcb1-dev \
                     libx11-xcb-dev libglu1-mesa-dev libxi-dev libxrandr-dev libxinerama-dev \
                     libxkbcommon-dev libxkbcommon-x11-dev \
                     libfontconfig1-dev libfreetype6-dev libssl-dev
```

- Qt 5.13.0
```
wget https://download.qt.io/archive/qt/5.13/5.13.0/qt-opensource-linux-x64-5.13.0.run
chmod +x qt-opensource-linux-x64-5.13.0.run
./qt-opensource-linux-x64-5.13.0.run
```

- Set the environment variables
```
export CGO_ENABLED=1
export QT_DIR=/home/<username>/Qt5.13.0
export QT_API=5.13.0
export QT_PKG_CONFIG=true
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export GO111MODULE=on
export QT_WEBKIT=true
export PATH=$PATH:$GOPATH/bin:$QT_DIR/5.13.0/gcc_64/bin
```

### qtsetup
- Clone therecipe/qt into your GOPATH/src directory, which avoids issues with module replacement.
```
git clone https://github.com/therecipe/qt.git $(go env GOPATH)/src/github.com/therecipe/qt
cd $(go env GOPATH)/src/github.com/therecipe/qt
go install ./cmd/qtsetup
go install ./cmd/qtmoc
go install ./cmd/qtminimal
go install ./cmd/qtrcc
go install ./cmd/qtdeploy
```

- Install the qtwebkit package
```
wget https://github.com/annulen/webkit/releases/download/qtwebkit-5.212.0-alpha2/qtwebkit-5.212.0_alpha2-qt59-linux-x64.tar.xz
tar -xvf qtwebkit-5.212.0_alpha2-qt59-linux-x64.tar.xz
```

- Run the following command to install the Qt dependencies.
```
# with the -test flag set to true, the installation will be tested
$(go env GOPATH)/bin/qtsetup test && $(go env GOPATH)/bin/qtsetup -test=false

# or
$(go env GOPATH)/bin/qtsetup
```