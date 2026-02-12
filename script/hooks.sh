#!/bin/sh

# プロジェクト内の全ての.goファイルを対象に go fmt, go vet, goimports を実行

# goimportsがインストールされているか確認
goimports_bin=$(command -v goimports)
if [ -z "$goimports_bin" ]; then
  echo "goimportsが見つかりません。インストールしてください: go install golang.org/x/tools/cmd/goimports@latest"
  exit 1
fi

# すべての .go ファイルを取得
go_files=$(git ls-files '*.go' -x 'client/*')

# go fmt
for file in $go_files; do
  go fmt "$file"
done

# go vet
go vet .
if [ $? -ne 0 ]; then
  echo "go vet でエラーが発生しました: $file"
  exit 1
fi

# goimports
for file in $go_files; do
  "$goimports_bin" -w "$file"
done

exit 0
