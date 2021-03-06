# Check-Rancher

It's a general purpose to monitore Rancher plateform with external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc.

The following program must be called by your monitoring tools and it return the status (nagios status normalization) with human messages and some times perfdatas.
This program called Rancher API to compute the state of your Rancher plateform.

You can use it to monitore the following bloc or your Rancher plateform:
- *Certificates expiration*: It's check that all your enabled certificates are stil  (not expired). You can use argument to specifiy the threshold in days before the certificates expire.
- *Hosts state in project*: It's check that all your hosts enabled in project are online (active). You need to specify the project name.
- *Stack state*: It's check that all services enabled on stack work as expected. You need to specify the stack name.


## Usage

### Global parameters

You need to set the Rancher API informations for all checks.

```sh
./check-rancher --rancher-url https://rancher.company.com --rancher-key my-key --rancher-secret my-secret ... 
```

You need to specify the following parameters:
- **--rancher-url**: It's your rancher URL. If you use `https` with enterprise PKI, you need to add your enterprise CA in the host certificate store.
- **--rancher-key**: It's the rancher API key that you have generated from Rancher UI (API -> Keys -> Add Account API Key)
- **--rancher-secret**: It's the secret associated to your API key.

### Check certificates expiration

You need to lauch the following command:

```sh
./check-rancher --rancher-url https://rancher.company.com --rancher-key my-key --rancher-secret my-secret check-certificates --warning-days 10
```

You need to specify the following parameters:
- **--warning-days**: The number of days before certificate expire. In this exemple, if certificate expire in 8 days, it will return warning.

This check follow this logic:
1. `OK` when there are no certificate
2. `OK` when all certificates are disabled
3. `OK` when all certificates not expired and not yet on threshold in days before to expire
4. `WARNING` when one of the certificates is on the threshold in days before to expire
5. `CRITICAL` when one of the certificates is expired

> All certificates that have problem is displayed on messages.

It's return the following perfdata:
- **nbActifCertificate**: the number of actif certificates
- **nbExpiredCertificate**: the number of expired certificates 


### Check hosts states on project

You need to lauch the following command:

```sh
./check-rancher --rancher-url https://rancher.company.com --rancher-key my-key --rancher-secret my-secret check-hosts-project --project-name production
```

You need to specify the following parameters:
- **--project-name**: The name of your Rancher project that you should to check the hosts state.

This check follow this logic:
1. `OK` when all hosts are online (active)
2. `OK` when project is disabled
3. `WARNING` when there are no host in project
4. `WARNING` when there all hosts are disabled
5. `CRITICAL` when one of the hosts is disconnected

> All hosts that have problem is displayed on messages.

It's return the following perfdata:
- **nbHost**: the number of hosts in project
- **nbHostOnline**: the number of online hosts

### Check the stack state

You need to lauch the following command:

```sh
./check-rancher --rancher-url https://rancher.company.com --rancher-key my-key --rancher-secret my-secret check-hosts-project check-stack --stack-name gitlab
```

You need to specify the following parameters:
- **--stack-name**: The name of your Rancher stack that you should to check the state.

This check follow this logic:
1. `OK` when stack is disabled
2. `OK` when all service that composed the stack are disabled
3. `OK` when all service that compose the stack are active (no instance problem)
4. `WARNING` when stack is empty (no service associated to it)
5. `WARNING` when one of the services is degraded (at least one instance online)
6. `CRITICAL` when one of services is down (all instance offline)

> All services that have problem is displayed on messages.

It's return the following perfdata:
- **nbService**: the number of service that composed the stack
- **nbInstance**: the number of instance that composed the stack
- **nbServiceOnline**: the number or services that online (even if degraded)

For each services:
- **SERVICE_NAME-nbInstanceOnline**: the number of instance that online for this service


## Contribute

Your PR are welcome, but please use develop branch and not the master branch.
You can open issues or enhance ;)