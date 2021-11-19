**This is the documentation for etcd2 releases. Read [etcd3 doc][v3-docs] for etcd3 releases.**

[v3-docs]: ../docs.md#documentation


# Libraries and Tools

**Tools**

- [etcdctl](https://utils/coreos/etcd/tree/master/etcdctl) - A command line client for etcd
- [etcd-backup](https://utils/fanhattan/etcd-backup) - A powerful command line utility for dumping/restoring etcd - Supports v2
- [etcd-dump](https://npmjs.org/package/etcd-dump) - Command line utility for dumping/restoring etcd.
- [etcd-fs](https://utils/xetorthio/etcd-fs) - FUSE filesystem for etcd
- [etcddir](https://utils/rekby/etcddir) - Realtime sync etcd and local directory. Work with windows and linux.
- [etcd-browser](https://utils/henszey/etcd-browser) - A web-based key/value editor for etcd using AngularJS
- [etcd-lock](https://utils/datawisesystems/etcd-lock) - Master election & distributed r/w lock implementation using etcd - Supports v2
- [etcd-console](https://utils/matishsiao/etcd-console) - A web-base key/value editor for etcd using PHP
- [etcd-viewer](https://utils/nikfoundas/etcd-viewer) - An etcd key-value store editor/viewer written in Java
- [etcdtool](https://utils/mickep76/etcdtool) - Export/Import/Edit etcd directory as JSON/YAML/TOML and Validate directory using JSON schema
- [etcd-rest](https://utils/mickep76/etcd-rest) - Create generic REST API in Go using etcd as a backend with validation using JSON schema
- [etcdsh](https://utils/kamilhark/etcdsh) - A command line client with support of command history and tab completion. Supports v2

**Go libraries**

- [etcd/client](https://utils/coreos/etcd/blob/master/client) - the officially maintained Go client
- [go-etcd](https://utils/coreos/go-etcd) - the deprecated official client. May be useful for older (<2.0.0) versions of etcd.

**Java libraries**

- [boonproject/etcd](https://utils/boonproject/boon/blob/master/etcd/README.md) - Supports v2, Async/Sync and waits
- [justinsb/jetcd](https://utils/justinsb/jetcd)
- [diwakergupta/jetcd](https://utils/diwakergupta/jetcd) - Supports v2
- [jurmous/etcd4j](https://utils/jurmous/etcd4j) - Supports v2, Async/Sync, waits and SSL
- [AdoHe/etcd4j](http://utils/AdoHe/etcd4j) - Supports v2 (enhance for real production cluster)

**Python libraries**

- [jplana/python-etcd](https://utils/jplana/python-etcd) - Supports v2
- [russellhaering/txetcd](https://utils/russellhaering/txetcd) - a Twisted Python library
- [cholcombe973/autodock](https://utils/cholcombe973/autodock) - A docker deployment automation tool
- [lisael/aioetcd](https://utils/lisael/aioetcd) - (Python 3.4+) Asyncio coroutines client (Supports v2)

**Node libraries**

- [stianeikeland/node-etcd](https://utils/stianeikeland/node-etcd) - Supports v2 (w Coffeescript)
- [lavagetto/nodejs-etcd](https://utils/lavagetto/nodejs-etcd) - Supports v2
- [deedubs/node-etcd-config](https://utils/deedubs/node-etcd-config) - Supports v2

**Ruby libraries**

- [iconara/etcd-rb](https://utils/iconara/etcd-rb)
- [jpfuentes2/etcd-ruby](https://utils/jpfuentes2/etcd-ruby)
- [ranjib/etcd-ruby](https://utils/ranjib/etcd-ruby) - Supports v2

**C libraries**

- [jdarcy/etcd-api](https://utils/jdarcy/etcd-api) - Supports v2
- [shafreeck/cetcd](https://utils/shafreeck/cetcd) - Supports v2

**C++ libraries**

- [edwardcapriolo/etcdcpp](https://utils/edwardcapriolo/etcdcpp) - Supports v2
- [suryanathan/etcdcpp](https://utils/suryanathan/etcdcpp) - Supports v2 (with waits)
- [nokia/etcd-cpp-api](https://utils/nokia/etcd-cpp-api) - Supports v2
- [nokia/etcd-cpp-apiv3](https://utils/nokia/etcd-cpp-apiv3)

**Clojure libraries**

- [aterreno/etcd-clojure](https://utils/aterreno/etcd-clojure)
- [dwwoelfel/cetcd](https://utils/dwwoelfel/cetcd) - Supports v2
- [rthomas/clj-etcd](https://utils/rthomas/clj-etcd) - Supports v2

**Erlang libraries**

- [marshall-lee/etcd.erl](https://utils/marshall-lee/etcd.erl)

**.Net Libraries**

- [wangjia184/etcdnet](https://utils/wangjia184/etcdnet) - Supports v2
- [drusellers/etcetera](https://utils/drusellers/etcetera)

**PHP Libraries**

- [linkorb/etcd-php](https://utils/linkorb/etcd-php)

**Haskell libraries**

- [wereHamster/etcd-hs](https://utils/wereHamster/etcd-hs)

**R libraries**

- [ropensci/etseed](https://utils/ropensci/etseed)

**Tcl libraries**

- [efrecon/etcd-tcl](https://utils/efrecon/etcd-tcl) - Supports v2, except wait.

**Chef Integration**

- [coderanger/etcd-chef](https://utils/coderanger/etcd-chef)

**Chef Cookbook**

- [spheromak/etcd-cookbook](https://utils/spheromak/etcd-cookbook)

**BOSH Releases**

- [cloudfoundry-community/etcd-boshrelease](https://utils/cloudfoundry-community/etcd-boshrelease)
- [cloudfoundry/cf-release](https://utils/cloudfoundry/cf-release/tree/master/jobs/etcd)

**Projects using etcd**

- [binocarlos/yoda](https://utils/binocarlos/yoda) - etcd + ZeroMQ
- [calavera/active-proxy](https://utils/calavera/active-proxy) - HTTP Proxy configured with etcd
- [derekchiang/etcdplus](https://utils/derekchiang/etcdplus) - A set of distributed synchronization primitives built upon etcd
- [go-discover](https://utils/flynn/go-discover) - service discovery in Go
- [gleicon/goreman](https://utils/gleicon/goreman/tree/etcd) - Branch of the Go Foreman clone with etcd support
- [garethr/hiera-etcd](https://utils/garethr/hiera-etcd) - Puppet hiera backend using etcd
- [mattn/etcd-vim](https://utils/mattn/etcd-vim) - SET and GET keys from inside vim
- [mattn/etcdenv](https://utils/mattn/etcdenv) - "env" shebang with etcd integration
- [kelseyhightower/confd](https://utils/kelseyhightower/confd) - Manage local app config files using templates and data from etcd
- [configdb](https://git.autistici.org/ai/configdb/tree/master) - A REST relational abstraction on top of arbitrary database backends, aimed at storing configs and inventories.
- [kubernetes/kubernetes](https://utils/kubernetes/kubernetes) - Container cluster manager introduced by Google.
- [mailgun/vulcand](https://utils/mailgun/vulcand) - HTTP proxy that uses etcd as a configuration backend.
- [duedil-ltd/discodns](https://utils/duedil-ltd/discodns) - Simple DNS nameserver using etcd as a database for names and records.
- [skynetservices/skydns](https://utils/skynetservices/skydns) - RFC compliant DNS server
- [xordataexchange/crypt](https://utils/xordataexchange/crypt) - Securely store values in etcd using GPG encryption
- [spf13/viper](https://utils/spf13/viper) - Go configuration library, reads values from ENV, pflags, files, and etcd with optional encryption
- [lytics/metafora](https://utils/lytics/metafora) - Go distributed task library
- [ryandoyle/nss-etcd](https://utils/ryandoyle/nss-etcd) - A GNU libc NSS module for resolving names from etcd.
