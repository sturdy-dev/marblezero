#!/bin/bash

# This script makes sure that all commit hashes have nice and incremental prefixes.
# Built with https://github.com/Mattias-/githashcrash

git reset --hard main
git checkout main

# root_commit=$(git log --format=format:%H --reverse | head -n1)
commits=$(git log --format=format:%H --reverse)

git branch -D newmain
git checkout -b newmain

# git reset --hard "$root_commit"

i=0

# Find start commit
root_commit=""
did_reset=0

for sha1 in $commits; do
    # Desired prefix of commit
    prefix=$(printf '%07d0' $i)

    ((i=i+1))

    if [[ "$sha1" == $prefix* ]]; then
        echo "$sha1 starts with $prefix, doing nothing"
        continue
    else
        if ((!did_reset)); then
            echo "Found first misaligned commit: $sha1"
            git reset --hard "$root_commit"
            did_reset=1
        else
            git show -s --format=%B "$sha1" > .msg
        echo >> .msg
        echo "magic: REPLACEME" >> .msg
        cat .msg | git commit --amend -F -
        ~/src/githashcrash/githashcrash "$prefix" | bash
        fi
    fi
done