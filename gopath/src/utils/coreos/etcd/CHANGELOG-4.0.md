

Previous change logs can be found at [CHANGELOG-3.x](https://utils/coreos/etcd/blob/master/CHANGELOG-3.x.md).


Planned breaking changes.


## v4.0.0 (TBD)

See [code changes](https://utils/coreos/etcd/compare/v3.3.0...v4.0.0) and [v4.0 upgrade guide](https://utils/coreos/etcd/blob/master/Documentation/upgrades/upgrade_4_0.md) for any breaking changes. **Again, before running upgrades from any previous release, please make sure to read change logs below and [v4.0 upgrade guide](https://utils/coreos/etcd/blob/master/Documentation/upgrades/upgrade_4_0.md).**

### Breaking Changes

- [Secure etcd by default](https://utils/coreos/etcd/issues/9475)?
- Change `/health` endpoint output.
  - Previously, `{"health":"true"}`.
  - Now, `{"health":true}`.
  - Breaks [Kubernetes `kubectl get componentstatuses` command](https://utils/kubernetes/kubernetes/issues/58240).
- Deprecate [`etcd --proxy*`](TODO) flags; **no more v2 proxy**.
- Deprecate [v2 storage backend](https://utils/coreos/etcd/issues/9232); **no more v2 store**.
  - v2 API is still supported via [v2 emulation](TODO).
- Deprecate [`etcdctl backup`](TODO) command.
- `clientv3.Client.KeepAlive(ctx context.Context, id LeaseID) (<-chan *LeaseKeepAliveResponse, error)` is now [`clientv4.Client.KeepAlive(ctx context.Context, id LeaseID) <-chan *LeaseKeepAliveResponse`](TODO).
  - Similar to `Watch`, [`KeepAlive` does not return errors](https://utils/coreos/etcd/issues/7488).
  - If there's an unknown server error, kill all open channels and create a new stream on the next `KeepAlive` call.
- Rename `utils/coreos/client` to `utils/coreos/clientv2`.

