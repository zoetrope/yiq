# jiq

It's [jid](https://github.com/simeji/jid) with [jq](https://stedolan.github.io/jq/).

You can drill down interactively by using [jq](https://stedolan.github.io/jq/) filtering queries.
jiq uses [jq](https://stedolan.github.io/jq/) internally, and it **requires** you to have "jq" in your PATH.

## Demo

![demo-jiq-main](https://github.com/simeji/jid/wiki/images/demo-jid-main-640.gif)

## Installation

* [Simply use "jiq" command](#simply-use-jiq-command)  
* [Build "jiq" command by yourself](#build-jiq-command-by-yourself)  

### Simply use "jiq" command

If you simply want to use `jiq` command, please download binary from below.

https://github.com/fiatjaf/jiq/releases

### Build "jiq" command by yourself

`go get github.com/fiatjaf/jiq/cmd/jiq`

## Usage

### Quick start

* [simple json example](#simple-json-example)  
* [another example](#another-example)
* [with curl](#with-curl)  

#### simple json example

execute the following command.

```
echo '{"aa":"2AA2","bb":{"aaa":[123,"cccc",[1,2]],"c":321}}'| jiq
```

then jiq will be running.

You can dig JSON data incrementally.

When you enter `.bb.aaa[2]`, you will see the following.

```
[Filter]> .bb.aaa[2]
[
  1,
  2
]
```

Then, you press Enter key and output `[1,2]` and exit.

#### another example

This json is used by [demo section](https://github.com/fiatjaf/jiq#demo).
```
echo '{"info":{"date":"2016-10-23","version":1.0},"users":[{"name":"fiatjaf","uri":"https://github.com/fiatjaf","id":1},{"name":"fiatjaf2","uri":"https://example.com/fiatjaf","id":2},{"name":"fiatjaf3","uri":"https://example.com/fiatjaf3","id":3}],"userCount":3}}'|jiq
```

#### with curl

Sample for using [RDAP](https://datatracker.ietf.org/wg/weirds/documents/) data.

```
curl -s http://rdg.afilias.info/rdap/domain/example.info | jiq
```

## Keymaps

|key|description|
|:-----------|:----------|
|`TAB` / `CTRL` + `I` |Show available items and choice them|
|`CTRL` + `W` |Delete from the cursor to the start of the word|
|`CTRL` + `F` / Right Arrow (:arrow_right:)|To the first character of the 'Filter'|
|`CTRL` + `B` / Left Arrow (:arrow_left:)|To the end of the 'Filter'|
|`CTRL` + `A`|To the first character of the 'Filter'|
|`CTRL` + `E`|To the end of the 'Filter'|
|`CTRL` + `J`|Scroll json buffer 1 line downwards|
|`CTRL` + `K`|Scroll json buffer 1 line upwards|
|`CTRL` + `L`|Change view mode whole json or keys (only object)|

### Option

-q : Print query (for jq)
