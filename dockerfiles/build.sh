#!/bin/bash

echo What image version of brian to publish?
echo -n "> "
read version

files=($(find -E . -type f -regex ".*.Dockerfile" -exec basename {} \;))

for file in ${files[*]}
do 
  lang=$(echo $file | cut -d . -f1)
  tag=robherley/brian-$lang:$version
  echo Building $lang image:
  docker build --build-arg VERSION=$version -t $tag -f $file .
  docker push $tag
done