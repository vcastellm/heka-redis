heka-redis
==========

Redis PubSub input plugin for [Mozilla Heka](http://hekad.readthedocs.org/)

See [Building *hekad* with External Plugins](http://hekad.readthedocs.org/en/latest/installing.html#build-include-externals)
for compiling in plugins.

Basically, you'll need to edit the cmake/plugin_loader.cmake file and add

    add_external_plugin(git https://github.com/victorcoder/heka-redis master)

And build heka

Debug
=====

Add this to your hekad.toml file to see log output of your subscriptions:

```toml
[PayloadEncoder]
append_newlines = false
prefix_ts = true
ts_format = "2006/01/02 3:04:05PM MST"

[debug]
type = "LogOutput"
message_matcher = "Type == 'redis_pub_sub'"
encoder = "PayloadEncoder"
```
