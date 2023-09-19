call "C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\VC\Auxiliary\Build\vcvars64.bat"
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\msvcp140.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\msvcp140_1.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\msvcp140_2.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\msvcp140_atomic_wait.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\msvcp140_codecvt_ids.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\vccorlib140.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\vcruntime140.dll" build\windows\
COPY "%VCToolsRedistDir%\x64\Microsoft.VC142.CRT\vcruntime140_1.dll" build\windows\
