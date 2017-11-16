# transline

Translate word or phrase using Yandex Services.

## Dictionary lookup

![dictionary text](https://i.imgur.com/PQKt62v.gif)

![dictionary json](https://i.imgur.com/9EqE1Uo.gif)

## Machinery Translation

![translation text](https://i.imgur.com/s89PFxD.gif)

![translation json](https://i.imgur.com/4XGWXD3.gif)

## Installation

```
go get github.com/kovetskiy/transline
```

Arch Linux users can install package from aur:

```
yaourt -Sy transline-git
```


## Options
- *-d --dictionary* - Use dictionary for translation.
- *-t --translator* - Use machinery translation.
- *-l --lang <lang>* - Translation direction [default: en-ru].
- *-s --synonyms <limit>* - Limit synonims. [default: 0]
- *-o --output <format>* - Output format. Can be text or json. [default: text]

