#/bin/sh -f
set -e

if [ ! -f "$HYPERSCAN_ROOT/lib/libhs.a" ]; then
	wget https://github.com/01org/hyperscan/archive/v$HYPERSCAN_VERSION.tar.gz -O /tmp/hyperscan.tar.gz
	mkdir -p /tmp/hyperscan
	tar -xzf /tmp/hyperscan.tar.gz -C /tmp/hyperscan --strip-components 1
	cd /tmp/hyperscan
	rm -rf tools

	cmake . -DCMAKE_BUILD_TYPE=RelWithDebInfo \
			-DBUILD_STATIC_AND_SHARED=on

	make
	make install
else
	echo "Using cached hyperscan v${HYPERSCAN_VERSION} @ ${HYPERSCAN_ROOT}.";
fi
