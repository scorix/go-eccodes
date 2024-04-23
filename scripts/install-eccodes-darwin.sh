#! /bin/sh

ECC_VERSION=2.35.0

curl -L https://confluence.ecmwf.int/download/attachments/45757960/eccodes-$ECC_VERSION-Source.tar.gz --output eccodes-Source.tar.gz
tar -xzf eccodes-Source.tar.gz

mkdir -p build-eccodes
cd build-eccodes
cmake -DCMAKE_INSTALL_PREFIX=$HOME/.eccodes \
    -DENABLE_FORTRAN=OFF \
    -DENABLE_MEMFS=ON \
    ../eccodes-$ECC_VERSION-Source &&
    make -j4 &&
    make install
