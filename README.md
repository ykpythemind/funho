# funho

## development

### requirements

* golang
* dep
* make

### quick start

```
$ docker-compose up -d // DB用にdockerコンテナを使う場合
$ dep ensure
$ make migrate
$ make seeds
$ make run
```

## mode

```
$ CONFIGOR_ENV=production make run
```

