mackerel-plugin-solrdih [![Build Status](https://travis-ci.org/supercaracal/mackerel-plugin-solrdih.svg?branch=master)](https://travis-ci.org/supercaracal/mackerel-plugin-solrdih)
=====================

Apache Solr DataImportHandler status metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-solrdih -core=<corename> [-url=<url>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.solrdih]
command = "/path/to/mackerel-plugin-solrdih -core=testcore"
```

## See also
* [Uploading Structured Data Store Data with the Data Import Handler](https://lucene.apache.org/solr/guide/8_1/uploading-structured-data-store-data-with-the-data-import-handler.html)
* [go-mackerel-pluginを利用してカスタムメトリックプラグインを作成する](https://mackerel.io/ja/docs/entry/advanced/go-mackerel-plugin)
* [mkr plugin installに対応したプラグインを作成する](https://mackerel.io/ja/docs/entry/advanced/make-plugin-corresponding-to-installer)
