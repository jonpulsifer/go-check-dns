# go-check-dns

This is a little whois thingy I made to track my domain name expirations in datadog.

My domains are hosted with [Google Cloud DNS](https://cloud.google.com/dns/) so it expects `GOOGLE_APPLICATION_CREDENTIALS` to be `/path/to/service-account.json` to fetch the managed domains.

It also expects `DATADOG_URL` to be set in ENV

Ref: https://cloud.google.com/compute/docs/access/service-accounts to learn about service accounts

![datadog](https://raw.githubusercontent.com/JonPulsifer/go-check-dns/master/go-check-dns.png)

## WARNING :fire:

Cron support will be dropped as soon as kubernetes cronjobs leave alpha on GKE. See https://github.com/kubernetes/features/issues/19 for details
## Configuration

Edit `crontab` with the GCP or script you want to run

### Kubernetes
@kelseyhightower has a great tutorial for service accounts: https://github.com/kelseyhightower/gke-service-accounts-tutorial

## Usage:

```sh
# exported variable
export DATADOG_URL=datadog.kube-public.svc.cluster.local:8125
./go-check-dns -project kubesec

# runtime variable
DATADOG_URL=127.0.0.1:8125 ./go-check-dns -project kubesec

# or use make for 127.0.0.1:8125 and specify a runtime gcp project
PROJECT=kubesec make run
```

## Credits:

Thanks to https://github.com/muchlearning/kubernetes-cron for cronspiration
