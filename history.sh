#!/bin/bash

# This script makes sure that all commit hashes have nice and incremental prefixes.
# Built with https://github.com/Mattias-/githashcrash

#git reset --hard main
#git checkout main

# All commits in our repository
commits=$(git log --format=format:%H --reverse)

git branch -D newmain
git checkout -b newmain

i=0

# Find start commit
prev_commit=""
did_reset=0

for sha1 in $commits; do
    # Desired prefix of commit
    prefix=$(printf '%07d0' $i)

    ((i=i+1))


    if [[ "$sha1" == $prefix* ]] && ((!did_reset)); then
        echo "$sha1 starts with $prefix, doing nothing"
        prev_commit="$sha1"
        continue
    else
        if ((!did_reset)); then
            echo "Found first misaligned commit=$sha1 parent=$prev_commit"
            git reset --hard "$prev_commit"
            did_reset=1   
        fi

        git cherry-pick "$sha1"

        git show -s --format=%B "$sha1" > .msg
        echo >> .msg
        echo "magic: REPLACEME" >> .msg
        cat .msg | git commit --amend -F -
        ~/src/githashcrash/githashcrash "$prefix" | bash
    fi
done

git branch -D main
git checkout -b main