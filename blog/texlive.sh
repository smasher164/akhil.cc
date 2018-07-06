#!/bin/sh
wget http://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz
tar -xzf install-tl-unx.tar.gz
cd install-tl-20*
./install-tl --profile=/texlive.profile
export PATH=/tmp/texlive/bin/x86_64-linux:$PATH
tlmgr install latex-bin dvisvgm standalone xkeyval #ifxetex ifpdf ifluatex
export PATH=/:$PATH
blog --conf /posts.toml