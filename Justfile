build:
  go build -o gitstats 

run *flags="":
  ./gitstats {{flags}}

build-and-run *flags="": build
  ./gitstats {{flags}}