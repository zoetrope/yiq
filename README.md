# jiq [![Mentioned in Awesome jq](https://awesome.re/mentioned-badge.svg)](https://github.com/fiatjaf/awesome-jq)

It's [jid](https://github.com/simeji/jid) with [jq](https://stedolan.github.io/jq/).

You can drill down interactively by using [jq](https://stedolan.github.io/jq/) filtering queries.

jiq uses [jq](https://stedolan.github.io/jq/) internally, and it **requires** you to have `jq` in your `PATH`.

If you prefer, there's an experimental, standalone, purely client-side web version of this on https://jq.alhur.es/jiq/.

## Demo

![screencast-repo.gif](screencast-repo.gif)

![screencast-packagejson.gif](screencast-packagejson.gif)

## Installation

```
go get github.com/fiatjaf/jiq/cmd/jiq
```

If you don't have `jq` installed, follow instructions at https://stedolan.github.io/jq/download/ and make sure to put it in your `PATH`.

## Usage

### Quick start

* [simple example](#simple-example)  
* [advanced usage examples](#advanced-usage-examples)
* [with curl](#with-curl)  

#### simple example

Execute the following command:

```
echo '{"aa":"2AA2","bb":{"aaa":[123,"cccc",[1,2]],"c":321}}'| jiq
```

Then jiq will be running. Now you can dig JSON data incrementally.

When you enter `.bb.aaa[2]`, you will see the following.

```
[Filter]> .bb.aaa[2]
[
  1,
  2
]
```

If you press Enter now it will output

```json
[
  1,
  2
]
```

and exit (if you want all the output in a single line you can either call `jiq -c` or pipe it into jq as `jiq | jq -c .`).

#### advanced usage examples

If you have ever used jq, you'll be familiar with these:

```
echo '{"economists": [{"id": 1, "name": "menger"}, {"id": 2, "name": "mises"}, {"name": "hayek", "id": 3}]}' | jiq
```

Now try writing `.economists | "\(.[0].name), \(.[1].name) and \(.[2].name) are economists."` or `[.economists.[].id]`, or even `.economists | map({key: "\(.id)", value: .name}) | from_entries`

#### with curl

Sample for using [RDAP](https://datatracker.ietf.org/wg/weirds/documents/) data.

```
curl -s http://rdg.afilias.info/rdap/domain/example.info | jiq
```

### command line arguments

`-q` : print the jq filter instead of the resulting filtered JSON to stdout (if you plan to use this with jq later)

Plus all the [arguments jq accepts](https://stedolan.github.io/jq/manual/#Invokingjq) -- they will affect both the JSON output inside jiq and the output that is printed to stdout (beware that some may cause bugs).

---

traffic analytics for this repo:

[![](https://ght.trackingco.de/fiatjaf/jiq)](https://ght.trackingco.de/)
