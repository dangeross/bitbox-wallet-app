:: Compiles the Qt5 app. Part of `make windows`, which also compiles/bundles the deps

call "C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\VC\Auxiliary\Build\vcvars64.bat"
cd build
qmake ..\BitBox.pro
nmake
