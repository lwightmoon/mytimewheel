# mytimewheel
usage
```
  w := mytimewheel.NewWheel(10, 1024)
  ticker := w.NewTicker(time.Second)
  for {
    <-ticker.C
    //do sth
   }
```
