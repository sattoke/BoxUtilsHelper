!include StrFunc.nsh
${Using:StrFunc} StrRep
!include ReplaceInFile.nsh
!include nsDialogs.nsh
!include LogicLib.nsh

Unicode true
LoadLanguageFile "${NSISDIR}\Contrib\Language files\Japanese.nlf"
Name "BoxUtilsHelper"
OutFile "Install.exe"
InstallDir "$LOCALAPPDATA\BoxUtilsHelper"
RequestExecutionLevel user
XPStyle on

Var Dialog
Var Label
Var Text
Var Text_State

Page license
Page custom nsDialogsPage nsDialogsPageLeave
Page instfiles

LicenseData LICENSE

Function .onInit
    StrCpy $Text_State "fgcoeikpcffiecgbanbaaoigcggkgihd"
FunctionEnd

Function nsDialogsPage
    nsDialogs::Create 1018
    Pop $Dialog

    ${If} $Dialog == error
        Abort
    ${EndIf}

    ${NSD_CreateLabel} 0 0 100% 12u "拡張機能のIDを入力してください。"
    Pop $Label

    ${NSD_CreateText} 0 13u 100% 12u $Text_State
    Pop $Text

    nsDialogs::Show
FunctionEnd

Function nsDialogsPageLeave
    ${NSD_GetText} $Text $Text_State
FunctionEnd

Section
    SetOutPath "$INSTDIR"
    File "boxutils-helper.exe"
    File "boxutils-helper.json"
    !insertmacro _ReplaceInFile "$INSTDIR\boxutils-helper.json" "__EXTENSION_ID__" "$Text_State"
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    WriteRegStr HKCU "Software\Google\Chrome\NativeMessagingHosts\jp.toke.boxutils_helper" "" "$INSTDIR\boxutils-helper.json"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\Uninstall.exe"
    Delete "$INSTDIR\boxutils-helper.exe"
    Delete "$INSTDIR\boxutils-helper.json"
    RMDir "$INSTDIR"
    DeleteRegKey HKCU "Software\Google\Chrome\NativeMessagingHosts\jp.toke.boxutils_helper"
SectionEnd
