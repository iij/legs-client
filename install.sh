#!/bin/sh

set -e

BIN_DIR=$PWD
BIN_PATH="${BIN_DIR}/legsc"

REPO="iij/legs-client"
REPO_URL="https://github.com/${REPO}/releases/latest"
API_URL="https://github.com/api/v3/repos/${REPO}/releases/latest"

OS=`uname -s | awk '{print tolower($0)}'`
ARCH=`uname -m | awk '{print tolower($0)}'`
GOARCH=''

PROXY=${http_proxy}
if [ -n "${https_proxy}" ]; then
	PROXY=${https_proxy}
fi

#########################

get_latest_release() {
	if [ -z "${PROXY}" ]; then
		curl --silent ${API_URL} |
			grep '"tag_name":' |
			sed -E 's/.*"([^"]+)".*/\1/'
	else
		curl --silent -x ${PROXY} ${API_URL} |
			grep '"tag_name":' |
			sed -E 's/.*"([^"]+)".*/\1/'
fi
}

get_goarch() {
	if echo $1 | egrep '(x86_64|amd64)' > /dev/null; then
		echo 'amd64'
	elif echo $1 | egrep '(i386|i686)' > /dev/null; then
		echo '386'
	elif echo $1 | egrep '(armv5|armv6|armv7)' > /dev/null; then
		echo 'arm'
	elif echo $1 | egrep 'armv8' > /dev/null; then
		echo 'arm64'
	fi
}

#########################

echo "install legsc from ${REPO_URL}"
echo ""

if ! echo ${OS} | egrep '(darwin|linux|netbsd)' > /dev/null; then
	echo "Sorry, unsupported os type by install script: ${OS}.\nPlease select binary in: ${REPO_URL}"
	exit 1
fi

GOARCH=`get_goarch ${ARCH}`

if [ -z "${GOARCH}" ];then
	echo "Sorry, unsupported architecture by install script: ${ARCH}.\nPlease select binary in: ${REPO_URL}"
	exit 1
fi

echo "OS: ${OS}"
echo "ARCH: ${GOARCH}"

LATEST_TAG=`get_latest_release`
BIN_URL="https://github.com/iij/legs-client/releases/download/${LATEST_TAG}/legsc_${OS}_${GOARCH}.tar.gz"

echo ""
echo "donwload binary from ${BIN_URL}"


if [ -z "${PROXY}" ]; then
	curl --silent -L ${BIN_URL} | tar -zx -f - -C ${BIN_DIR}
else
	curl --silent -x ${PROXY} -L ${BIN_URL} | tar -zx -f - -C ${BIN_DIR}
fi

mv "${BIN_PATH}_${OS}_${GOARCH}" ${BIN_PATH}

echo ""
echo "save to ${BIN_PATH}"
echo "chmod 755 ${BIN_PATH}"

chmod 755 ${BIN_PATH}

echo ""
echo 'create config dir(~/$XDG_CONFIG_HOME(default: .config)/legsc/)'
if [ -z "${XDG_CONFIG_HOME}" ]; then
	mkdir -p ~/.config/legsc
else
	mkdir -p ~/$XDG_CONFIG_HOME/legsc
fi


echo ""
echo "If install was successful, you can see the version info."
echo "--------------------------------"
${BIN_PATH} version
echo "--------------------------------"

echo ""
echo "Use \"./legsc help\" for more information about legsc command."
