# cbr2mobi

Convert cbr files to mobi keeping the folder structure.

# Requirements

- https://calibre-ebook.com/

# Usage

Usage of cbr2mobi:

    $ cbr2mobi -h
        -i string
            directory of origin
        -o string
            directory of destination
        -v	verbose output

It will look for every file inside the folder specified recursively, and output the converted files to the destination folder, keeping the same folder structure.

    $ cbr2mobi -i /my/folder/of/origin -o /my/folder/of/destination -v

Destination folder will be equal to origin by default, if not specified.

    $ cbr2mobi -i /my/folder/of/origin
