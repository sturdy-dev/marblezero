#!/bin/bash

# This script makes sure that all commit hashes have nice and incremental prefixes.
# Built with https://github.com/Mattias-/githashcrash

git reset --hard main
git checkout main

root_commit=$(git log --format=format:%H --reverse | head -n1)
commits=$(git log --format=format:%H --reverse)

git branch -D newmain
git checkout -b newmain

git reset --hard "$root_commit"

i=0

for sha1 in $commits; do
    # Desired prefix of commit
    prefix=$(printf '%07d0' $i)    

    if ((i!=0)); then
        git cherry-pick "$sha1"
    fi

    if [[ "$sha1" == $prefix* ]]; then
        echo "$sha1 starts with $prefix"
    else
        git show -s --format=%B "$sha1" > .msg
        echo >> .msg
        echo "magic: REPLACEME" >> .msg
        cat .msg | git commit --amend -F -
        ~/src/githashcrash/githashcrash "$prefix" | bash
    fi

    ((i=i+1))
done