cmake_minimum_required(VERSION 3.10)

project (scan_test)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED True)

find_package (benchmark REQUIRED)
find_package(PkgConfig REQUIRED)

pkg_check_modules (HYPERSCAN REQUIRED IMPORTED_TARGET libhs)

add_executable (scan_test scan_test.cpp)

target_include_directories (scan_test PRIVATE ${HYPERSCAN_INCLUDE_DIRS})
target_link_libraries (scan_test benchmark::benchmark PkgConfig::HYPERSCAN)
