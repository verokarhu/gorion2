gorion2
=======

Master of Orion 2 rewrite

At the moment the only usable thing here is lbxdumper which can be installed by:

```
go get github.com/verokarhu/gorion2/cmd/lbxdumper
```

It supports video, image and audio lbx files. Using it is simple:

```
lbxdumper
  -a=false: assume audio content
  -cpuprofile="": write cpu profile to file
  -dir="dumpdir": directory where the dumped files go
  -game="disc": path to directory containing moo2 install or the contents of the game disc
  -i=false: assume image content
  -lbx="": name of lbx file
  -pal="all": name of palette to use, list lists the alternatives
  -v=false: assume video content
```