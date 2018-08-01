#/bin/sh -f
set -e

# things to do for travis-ci in the before_install section

if [[ $TRAVIS_OS_NAME == 'osx' ]]; then
	brew update
	brew outdated cmake || brew upgrade cmake
	brew outdated boost || brew upgrade boost
	brew install ragel tree
fi
