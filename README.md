# kubetrack

Tracking Kubernetes resources and events as a daemon

## Movivation

In production, you may need to track the difference of you kubernetes resources or to find the evidence of modifying resources.

Kubetrack watches every resources and events you cares about, then extracts what you want into databases or print logs.

By default, kubernetes events will be stored only 1 hour, kubetrack will store the events into databases and you can set TTL you want.

## Get started

Let's install using helm.

```bash
cd kubetrack
helm upgrade --install --create-namespace --namespace kubetrack kubetrack deploy/chart/kubetrack
```

## Configuration

The example configuration is put here: `conf/config.yaml`

```yaml
apiVersion: kubetrack.io/v1
kind: KubeTrackConfiguration
rules:
  - apiVersion: "v1"
    # the resource kind you need to watch
    kind: Pod

    # namespaces you need to watch
    namespaces: ["default"]

    # excluded namespaces, wildcard is supported here
    excludedNamespaces: ["dontcare", "kube*"]

    # fields you really care about, these fields will be stored seperate
    #   supported types are jsonpath, go-template, builtin
    #   - for jsonpath, the syntax you may refer to https://kubernetes.io/docs/reference/kubectl/jsonpath/
    #   - builtin functions are listed here:
    #       - PodStatus
    #       - PodStatusWithRestartCount
    #       - NodeStatus
    #       - FindNodeRoles
    careFields:
      - name: deletionTimestamp
        type: jsonpath
        expr: .metadata.deletionTimestamp
      - name: phase
        type: jsonpath
        expr: .status.phase
      - name: status
        type: builtin
        expr: PodStatusWithRestartCount
      - name: nodeName
        type: jsonpath
        expr: .spec.nodeName

    # whether you need to save the full object, save the comparasion(diff) or jsonpatch on create/update/delete
    onCreate:
      saveFullObject: true
    onUpdate:
      saveFullObject: false
      saveCmp: true
      saveJsonPatch: true
    onDelete:
      saveFullObject: true
  - apiVersion: "v1"
    kind: Node
    careFields:
      - name: status
        type: builtin
        expr: NodeStatus
      - name: role
        type: builtin
        expr: FindNodeRoles
      - name: addresses
        type: jsonpath
        expr: .status.addresses
      - name: podCIDR
        type: jsonpath
        expr: .spec.podCIDR
      - name: taints
        type: jsonpath
        expr: .spec.taints
    onCreate:
      saveFullObject: true
    onUpdate:
      saveFullObject: false
      saveCmp: true
      saveJsonPatch: true
    onDelete:
      saveFullObject: true

# what namespaces to watch
events:
  namespaces: [] # watch all namespaces
  excludedNamespaces: []

# save the output to one or multiple the databases
output:
  - log:
      printDiff: true
  - postgres:
      dsn: host=127.0.0.1 user=postgres password=password dbname=kubetrack port=5432 sslmode=disable connect_timeout=5
      ttlDays: 1
  - mysql:
      dsn: "root:password@tcp(127.0.0.1:3306)/kubetrack?charset=utf8mb4&parseTime=True&loc=Local"
      ttlDays: 1
```

## Useful SQLs

Show latest 10 records

```sql
SELECT id, event_time, source, event_type,
  concat(api_version, '/', kind, ',', namespace, '/', name) object_ref,
  uid, message, fields
FROM events
ORDER BY id DESC
LIMIT 10\G
```

Show most produced events group by kind, namespace and source

```sql
SELECT api_version, kind, namespace, source, count(*) cnt
FROM events
GROUP BY api_version, kind, namespace, source
ORDER BY cnt DESC
```

## License

MIT
