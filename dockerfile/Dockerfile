# syntax=docker/dockerfile:1.4

# 上記で新しいsyntaxであることをビルダーに伝える。
# 新しい構文を使うとき、
# なぜかなくても動いたり動かったりする環境があってややこしいので
# とりあえず書く。

FROM golang:1.22.3-bookworm AS builder

ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG GOPATH=/go
ARG CGO_ENABLED=0
ARG MAIN_PKG_PATH=.

# WORKDIRの決め方やビルドしたバイナリの置き場所はこれがいいよという自信がない。
# 必要に応じて変えてください。
WORKDIR /usr/local/container-bin/src
# git-lfsの有無でgit fetch結果が異なり、sum照合エラーになることがある。
# Private go moduleをdirect modeでgo getするならば、すべての環境に入れておくほうが安全。
# apt-getでバージョン指定をするとすぐに古いパッケージが消えるのでバージョンは固定しない。
# バージョンを固定したい場合はdebファイルを保存して
# そこからインストールしたり、ソースからビルドする。
RUN --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt\
    apt-get update && apt-get install -yqq --no-install-recommends git-lfs
# 先にgo mod downloadを実行する
# buildkitでマウントするキャッシュ以外に変更が起きない。
# (/root/.cacheと/root/.config/goにマウントされるのでディレクトリは作成される)
# Dockerのimage layerとしてキャッシュするというより、
# コマンドの失敗する点を切り分けてエラーを見やすくする意図がある。
COPY go.mod go.sum ./
RUN --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt\
    --mount=type=secret,id=.netrc,target=/root/.netrc\
    --mount=type=secret,id=goenv,target=/root/.config/go/env\
    --mount=type=cache,target=/go\
    --mount=type=cache,target=/root/.cache/go-build\
    go mod download
# COPY . .をしてしまうとbuildkitの遅延ファイル要求の利点がすっ飛びますが、全部送らざるを得ない
# ソース以外のコンテンツがいろいろ含まれる場合は、`.dockerignore`などをちきんと整備してください。
# https://docs.docker.com/build/building/context/#dockerignore-files
COPY . .
RUN --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt\
    --mount=type=secret,id=.netrc,target=/root/.netrc\
    --mount=type=secret,id=goenv,target=/root/.config/go/env\
    --mount=type=cache,target=/go\
    --mount=type=cache,target=/root/.cache/go-build\
    go build -o ../bin ${MAIN_PKG_PATH}

# distrolessはtagの中身が入れ替わるので再現性を優先するならsha256で指定したほうがよい
FROM gcr.io/distroless/static-debian12@sha256:41972110a1c1a5c0b6adb283e8aa092c43c31f7c5d79b8656fbffff2c3e61f05

COPY --from=builder /usr/local/container-bin/bin /usr/local/container-bin/

ENTRYPOINT [ "/usr/local/container-bin/bin" ]