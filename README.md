m3usplit
========
My first attempt at writing a go language program.  It splits large m3u files
into smaller m3u files.

For example, if I had an m3u file containing the entire mp3 discography of deadmau5 and I wanted to split it into multiple m3u files that could be burned onto a 700M CDR:

m3usplit deadmau5_discog.m3u 700


Build
```
go build m3usplit.go
```

Run
```
./m3usplit
```


