# CV source

Build from the repo root (Computer Modern via pdflatex; needs `geometry`, `tabularx`, `hyperref`):

    nix shell nixpkgs#texliveMedium -c make -C static cv

This writes `static/root/cv.pdf`, which `make -C static deploy` serves at www.static.akhil.cc/cv.pdf.

Adding a job: copy an existing `\entry{Title}{Location}{Organisation}{Dates}` block in the
EXPERIENCE section (newest first) and edit it. Description text goes right after the `\entry`
line; an optional `\begin{bullets} \item ... \end{bullets}` follows. Projects and skills are
`\row{Label}{text}` lines inside `\begin{labeled}`; `\gh{name}` links a `github.com/smasher164` repo.
