# github.com/cbh34680/fmtsshconf

## Description

OpenSSH の config ファイルをフォーマットする。  
  
-config 引数が指定されない場合は、Windows では %USERPROFILE%\.ssh\config  
Linux では $HOME/.ssh/config が更新対象として採用されます。  

実行すると、対象ファイルにある重複する Host ブロックを一つにまとめますが  
同一ブロックに複数の重複するパラメータ (ex. "HostName") が存在する場合は  
最後に出現したものが有効になります。  

この動作は ssh コマンドが認識するパラメータが先頭のものであることと  
異なるため、注意が必要です。  

また、"#" から始まるコメントは全て削除されます。  


## Usage

```dos
> fmtsshconf -config C:\WORK\.ssh\config ... 更新対象を直接指定
> fmtsshconf -confirm=false ... 更新の確認をしない
> fmtsshconf -delhost host-of-delete.local ... 引数の Host ブロックを削除
```

