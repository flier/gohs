#/bin/sh -f
set -e

# things to do for travis-ci in the before_install section

if [[ $TRAVIS_OS_NAME == 'osx' ]]; then
	brew update
	brew outdated cmake || brew upgrade cmake
	brew outdated boost || brew upgrade boost
	brew install ragel
	brew install tree
else
	mkdir -p $HOME/bin

	ln -s /usr/bin/g++-4.8 $HOME/bin/g++
	ln -s /usr/bin/gcc-4.8 $HOME/bin/gcc
	ln -s /usr/bin/gcov-4.8 $HOME/bin/gcov

    export PATH=$HOME/bin:$PATH

	if [ ! -f "$BOOST_ROOT/lib/libboost_graph.a" ]; then
		wget http://downloads.sourceforge.net/project/boost/boost/1.$BOOST_VERSION_MINOR.0/boost_1_$BOOST_VERSION_MINOR\_0.tar.gz -O /tmp/boost.tar.gz
		mkdir -p /tmp/boost
		tar -xzf /tmp/boost.tar.gz -C /tmp/boost --strip-components 1
		cd /tmp/boost
		./bootstrap.sh
		./b2 -q -d=0 install -j 2 --prefix=$BOOST_ROOT link=static
	else
  		echo "Using cached boost v1.${BOOST_VERSION_MINOR}_0 @ ${BOOST_ROOT}.";
  	fi
fi
