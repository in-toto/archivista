#! /bin/sh

script_dir=$(dirname "$0")
root_dir=$(dirname "$script_dir")



cd "$root_dir" || exit


web_dir="$PWD/../web"

cd "$web_dir" || exit
#check if public folder is empty
if [ ! "$(ls -A public)" ]; then
    echo "public folder is empty, building web"
    npm ci
    npm run build
fi
