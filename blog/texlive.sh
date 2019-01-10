#!/bin/sh
wget http://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz
tar -xzf install-tl-unx.tar.gz
cd install-tl-20*
./install-tl --profile=/texlive.profile
export PATH=/tmp/texlive/bin/x86_64-linux:$PATH
# tlmgr install latex-tools
tlmgr install l3kernel l3packages l3experimental
tlmgr install latex-bin dvisvgm standalone xkeyval algorithms algorithmicx float amsmath siunitx #ifxetex ifpdf ifluatex
export PATH=/:$PATH
