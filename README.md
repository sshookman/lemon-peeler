# lemon-peeler
Site scraping tool for downloading files in bulk

---

Install
```
go get github.com/sshookman/lemon-peeler
```

Example Command - Downloads all of the [MagPi Magazines](https://www.raspberrypi.org/magpi/issues/) in PDF format
```
lemon-peeler -u https://www.raspberrypi.org/magpi/issues/ -s pdf -d
```
