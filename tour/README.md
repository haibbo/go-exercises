# Exercises in tour.golang.org

访问 https://tour.golang.org/welcome/1 获取tour完整的信息, 如果你无法访问这个网站, 可以通过以下方式创建本地版的tour.

```
git clone https://github.com/golang/tour.git $GOPATH/src/golang.org/x/tour
cd $GOPATH/src                                                     
go build golang.org/x/tour/gotour
mv gotour $GOPATH/bin
```

运行 gotour, 它会打开本地浏览器, 访问tour
```
~/g/src ❯❯❯ gotour                                                  
2018/03/17 23:43:32 Serving content from /Users/haibbo/go/src/golang.org/x/tour
2018/03/17 23:43:32 A browser window should open. If not, please visit http://127.0.0.1:3999
2018/03/17 23:43:35 accepting connection from: 127.0.0.1:50579
```

