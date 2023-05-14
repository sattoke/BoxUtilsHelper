# BoxUtilsHelper
BoxUtilsHelper はChrome拡張機能 BoxUtils ( [chrome ウェブストア](https://chrome.google.com/webstore/detail/boxutils/gagpkhipmdbjnflmcfjjchoielldogmm) 、 [GitHub](https://github.com/sattoke/BoxUtils) ) と連動して動作する補助ソフトである。
BoxUtilsから直接ローカルのOS上でフォルダやファイルを開く機能を使うときに必要となる。

# インストーラーの入手方法
GitHubからインストーラーのバイナリファイルを入手する方法と、GitHubにあるコードからインストーラーを作成する方法の2種類がある。
後者は開発者向けの方法のため、通常はバイナリファイルを入手すること(インストールが簡単)。

## GitHubからインストーラーのバイナリファイルを入手する場合 (一般利用者向け：通常はこちらを利用すること)
[GitHubのBoxUtilsHelperのリリースページ](https://github.com/sattoke/BoxUtilsHelper/releases) から最新のリリースに含まれる `Install-x.y.z.x86_64.exe` をダウンロードする。 ( `x.y.z` はバージョン番号 )

## GitHubにあるコードからインストーラーを作成する場合 (開発者向け)

- 前提条件
  - WSLを利用していること (付属のMakefileを使う場合)
  - NSISがインストールされていること
  - goがインストールされていること

1. https://github.com/sattoke/BoxUtilsHelper をgitでcloneするか、当該URLの「Code」ボタンをクリックすると出てくる「Download ZIP」でダウンロードし、ZIPを適当なところに展開する。

    ```sh
    $ git clone git@github.com:sattoke/BoxUtilsHelper.git
    ```

1. WSLのシェルから、リポジトリのトップディレクトリで `make` を実行するとインストーラー (`Install-x.y.z.x86_64.exe`) が作成される。

    ```sh
    $ cd BoxUtilsHelper
    $ make
    ```

# インストール方法
BoxUtilsHelperをインストールする前にChrome拡張機能 [BoxUtils](https://github.com/sattoke/BoxUtils) をインストールしておくこと。
このBoxUtilsのインストール方法に応じて、BoxUtilsHelperのインストール手順が異なるため、
「BoxUtilsをchrome ウェブストアからインストールした場合」と「BoxUtilsをchrome ウェブストア以外(GitHub)からインストールした場合」の適切な方に従ってインストールすること。
いずれの方法でインストールした場合でも、インストール先は `%LOCALAPPDATA%\BoxUtilsHelper` となる。

## BoxUtilsをchrome ウェブストアからインストールした場合 (一般利用者向け：通常はこちらを利用すること)
1. インストーラー (`Install-x.y.z.x86_64.exe`) を実行
1. 拡張機能のIDの入力を求められる画面ではデフォルトの値 (`gagpkhipmdbjnflmcfjjchoielldogmm`) のままとし、後はインストーラの指示通りに進めればよい。

## BoxUtilsをchrome ウェブストア以外(GitHub)からインストールした場合 (開発者向け)
1. Chromeのアドレスバーに `chrome://extensions/` と入力するか、Chromeのメニュー（ケバブメニュー）から、「設定」→「拡張機能」と選択することで拡張機能の管理画面を開く
1. 拡張機能一覧からBoxUtilsを探し、そこの `ID` 欄に記載されている拡張機能のID(英字32桁)をメモする。
1. BoxUtilsHelperのインストーラー (`Install-x.y.z.x86_64.exe`) を実行する。
1. 拡張機能のIDの入力を求められるためメモした拡張機能のIDを入力して、後はインストーラの指示通りに進める。


# 使用方法
BoxUtilsのポップアップメニューからフォルダやファイルを開くアイコンをクリックするとBoxUtilsHelperが呼び出されてBoxDrive経由でフォルダやファイルが開かれる(BoxUtilsHelperの実行ファイルなどを直接実行することはない)。


# アンインストール方法
アンインストールには下記のいずれかを実施する。

- BoxUtilsHelperのインストール先( `%LOCALAPPDATA%\BoxUtilsHelper` )にある Uninstall.exe を実行する。
- Windowsで「設定」→「アプリ」→「インストールされているアプリ」と辿り (もしくは「ファイル名を指定して実行」で `ms-settings:appsfeatures` を実行) 、BoxUtilsHelperをアンインストール
- Windowsで「コントロールパネル」→「プログラムのアンインストール」と辿り (もしくは「ファイル名を指定して実行」で `appwiz.cpl` を実行) 、BoxUtilsHelperからアンインストール
(
