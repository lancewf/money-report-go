pkg_name=moeney-report-go
pkg_description="Money report with a go backend"
pkg_origin=lancewf
pkg_version=$(cat "$PLAN_CONTEXT/../VERSION")
pkg_maintainer="Lance Finfrock <lancewf@gmail.com>"
pkg_license=('UNLICENSED')
pkg_upstream_url="http://github.com/lancewf/money-report-go"
pkg_build_deps=()
pkg_exports=(
  [port]=service.port
  [host]=service.host
)
pkg_exposes=(port)
pkg_binds_optional=(
  [elasticsearch]="http-port"
)
pkg_deps=(core/musl)
pkg_bin_dirs=(bin)
pkg_scaffolding=afiune/scaffolding-go
scaffolding_go_base_path=github.com/chef
scaffolding_go_build_deps=(
  github.com/golang/dep/cmd/dep
)

scaffolding_go_get_with_flags() {
  local deps
  deps=($pkg_source ${scaffolding_go_build_deps[@]})
  build_line "Downloading Go build dependencies"
  if [[ "${#deps[@]}" -gt 0 ]] ; then
    for dependency in "${deps[@]}" ; do
      go get --ldflags "${GO_LDFLAGS}" "$(_sanitize_pkg_source "$dependency")"
    done
  fi
}

do_download(){
  # Since we are building static lib
  build_line "Setting up CC/LDFLAGS for static linking"
  export CC=$(pkg_path_for core/musl)/bin/musl-gcc
  export GO_LDFLAGS="-linkmode external -extldflags '-static'"
  export GO_IMPORT_PATH="${scaffolding_go_base_path}/${pkg_name}"
  export PATH=$GOPATH/bin:$PATH
  build_line "Overriding scaffolding do_default_download"
  scaffolding_go_get_with_flags
  build_line "Installing go dependencies"
  pushd $scaffolding_go_pkg_path >/dev/null
  dep ensure
  popd >/dev/null
}

do_prepare() {
  export GIT_SHA=$(git rev-parse --short HEAD)
  export BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
  LINKER_FLAGS=" -X $GO_IMPORT_PATH/config.VERSION=$pkg_version"
  LINKER_FLAGS="$LINKER_FLAGS -X $GO_IMPORT_PATH/config.SHA=$GIT_SHA"
  LINKER_FLAGS="$LINKER_FLAGS -X $GO_IMPORT_PATH/config.BUILD_TIME=$BUILD_TIME"
  export LINKER_FLAGS
}

do_build() {
  build_line "Overriding Build process"
  pushd "$scaffolding_go_pkg_path" >/dev/null
  go build --ldflags "${GO_LDFLAGS} ${LINKER_FLAGS}"
  popd >/dev/null
}

do_install() {
  build_line "Overriding Install process"
  pushd "$scaffolding_go_pkg_path" >/dev/null
  go install --ldflags "${GO_LDFLAGS} ${LINKER_FLAGS}"
  popd >/dev/null
  cp -r "${scaffolding_go_gopath:?}/bin/${pkg_name}" "${pkg_prefix}/bin"
}
