#!/bin/bash
rm -rf flat
mkdir flat
for f in file_{a,b,c,d,e,f}.txt; do echo $f > ./flat/$f; done;

rm -rf tree
mkdir tree
for f in file_{a,b,c}.txt; do echo $f > ./tree/$f; done;
mkdir tree/subtree
for f in file_{d,e,f}.txt; do echo $f > ./tree/subtree/$f; done;

rm -rf zip
mkdir zip

pushd flat
zip -rm ../zip/flat.zip *
popd

pushd tree
zip -rm ../zip/tree.zip *
popd
