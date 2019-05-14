#/bin/sh -f
set -e

wget https://github.com/intel/hyperscan/archive/v$HYPERSCAN_VERSION.tar.gz -O /tmp/hyperscan.tar.gz
mkdir -p /tmp/hyperscan
tar -xzf /tmp/hyperscan.tar.gz -C /tmp/hyperscan --strip-components 1
cd /tmp/hyperscan
rm -rf tools

cmake . -DCMAKE_BUILD_TYPE=RelWithDebInfo \
		-DBUILD_STATIC_AND_SHARED=on \
		-DCMAKE_POSITION_INDEPENDENT_CODE=on

make
make install
