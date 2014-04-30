heka-redis
==========

Redis PubSub input plugin for [Mozilla Heka](http://hekad.readthedocs.org/)

See [Building *hekad* with External Plugins](http://hekad.readthedocs.org/en/latest/installing.html#build-include-externals)
for compiling in plugins.

Basically, you'll need to edit the cmake/plugin_loader.cmake file and add

    add_external_plugin(git https://github.com/victorcoder/heka-redis master)

And build heka
