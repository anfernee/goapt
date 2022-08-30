#!/usr/bin/bash

FILE="random.txt"

echo MD5
md5sum ${FILE}

echo SHA1
sha1sum ${FILE}

echo SHA256
sha256sum ${FILE}