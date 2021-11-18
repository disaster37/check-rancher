[![build](https://github.com/disaster37/check-rancher/actions/workflows/workflow.yaml/badge.svg)](https://github.com/disaster37/check-rancher/actions/workflows/workflow.yaml)
[![GoDoc](https://godoc.org/github.com/disaster37/check-rancher/v2?status.svg)](http://godoc.org/github.com/disaster37/check-rancher/v2)
[![codecov](https://codecov.io/gh/disaster37/check-rancher/branch/2.x/graph/badge.svg)](https://codecov.io/gh/disaster37/check-rancher/branch/2.x)


# Check-Rancher

It's a general purpose to monitore Rancher plateform with external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc.

The following program must be called by your monitoring tools and it return the status (nagios status normalization) with human messages and some times perfdatas.
This program called Rancher API to compute the state of your Rancher plateform.

You can use it to monitore the following bloc or your Rancher2 plateform:
- snapshot: ETCD backup state

## Contribute

Your PR are welcome, but please use 2.x branch and not the master branch.
You can open issues or enhance ;)

## Usage

### Global parameters

You need to set the Rancher API informations for all checks.

```sh
./check-rancher --url https://rancher.company.com --access-key my-key --secret-key ... 
```

You need to specify the following parameters:
- **--url**: It's your rancher URL. If you use `https` with enterprise PKI, you need to add your enterprise CA in the host certificate store or disable the certificate check.
- **--access-key**: It's the rancher API key that you have generated from Rancher UI
- **--secret-key**: It's the secret associated to your API key.

### Check snapshot (ETCD backup)

You can check that the last backup is OK and not older than.

You need to lauch the following command:

```sh
./check-rancher --url https://rancher.company.com --access-key my-key --secret-key check-snapshot --cluster-name my-k8s-cluster max-older-than 24h
```

You need to specify the following parameters:
- **--cluster-name**: The kubernetes cluster managed by rancher where you should to check ETCD backups.
- **--max-older-then**: The maxiumum duration the last backup can be not older than. Default to `24h`.

This check follow this logic:
- `OK` when last backup is OK and not older than
- `CRITICAL` when last backup is KO, or older than or no backup


It's return the following perfdata:
- **nbBackupOk**: the number of backup on active state
- **nbBackupFailed**: the number of backup failed 


