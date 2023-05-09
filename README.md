# BoxUtilsHelper
BoxUtilsHelper はChrome拡張機能 BoxUtils ( [chrome ウェブストア](https://chrome.google.com/webstore/detail/boxutils/gagpkhipmdbjnflmcfjjchoielldogmm) 、 [GitHub](https://github.com/sattoke/BoxUtils) ) と連動して動作する補助ソフトである。
BoxUtilsから直接ローカルのOS上でフォルダやファイルを開く機能を使うときに必要となる。

# インストーラーの作成方法
## 前提条件
- WSLを利用していること (付属のMakefileを使う場合)
- NSISがインストールされていること
- goがインストールされていること

WSLのシェルから `make` を実行するとInstall.nsiと同じフォルダにインストーラー (Install.exe) が作成される。

```
$ make
```

# インストール方法
## BoxUtilsをchrome ウェブストアからインストールした場合
1. インストーラー (Install.exe) を実行
1. 拡張機能のIDの入力を求められる画面ではデフォルトの値 ( `gagpkhipmdbjnflmcfjjchoielldogmm` ) のまま、後はインストーラの指示通りに進めればよい。

## BoxUtilsをchrome ウェブストア以外(GitHub)からインストールした場合

1. 予めChrome拡張機能 [BoxUtils](https://github.com/sattoke/BoxUtils) をインストールしておく
1. Chromeのアドレスバーに `chrome://extensions/` と入力するか、Chromeのメニュー（ケバブメニュー）から、「設定」→「拡張機能」と選択することで拡張機能の管理画面を開く
1. 拡張機能一覧からBoxUtilsを探し、そこの `ID` 欄に記載されている拡張機能のID(英字32桁)をメモする。
1. BoxUtilsHelperのインストーラー (Install.exe) を実行する。
1. 拡張機能のIDの入力を求められるためメモした拡張機能のIDを入力して、後はインストーラの指示通りに進める。

インストール先は `%LOCALAPPDATA%\BoxUtilsHelper` となる。

# 使用方法
BoxUtilsのポップアップメニューからフォルダやファイルを開くアイコンをクリックするとBoxUtilsHelperが呼び出されてBoxDrive経由でフォルダやファイルが開かれる(BoxUtilsHelperの実行ファイルなどを直接実行することはない)。


# アンインストール方法
BoxUtilsHelperのインストール先にある Uninstall.exe を実行する。
