cmake_minimum_required(VERSION 3.10)

project (scan_test)

set (CMAKE_CXX_STANDARD 11)
set (CMAKE_CXX_STANDARD_REQUIRED True)

option (USE_RE2 "benchmark re2" ON)

find_package (benchmark REQUIRED)
find_package (PkgConfig REQUIRED)

pkg_check_modules (HYPERSCAN REQUIRED IMPORTED_TARGET libhs)

add_executable (scan_test scan_test.cpp)

if (USE_RE2)
    pkg_check_modules (RE2 REQUIRED IMPORTED_TARGET re2)

    if (RE2_FOUND)
        target_include_directories (scan_test PRIVATE ${RE2_INCLUDE_DIRS})
        target_link_libraries (scan_test PkgConfig::RE2)
        set (HAVE_RE2 1)
    else ()
        message (FATAL_ERROR "re2 not found")
    endif ()
endif ()

configure_file (scan_test.h.in scan_test.h @ONLY)
include_directories(${CMAKE_CURRENT_BINARY_DIR})

target_include_directories (scan_test PRIVATE ${HYPERSCAN_INCLUDE_DIRS})
target_link_libraries (scan_test benchmark::benchmark PkgConfig::HYPERSCAN)
