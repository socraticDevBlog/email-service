# post messages to a third party API on a schedule

On Unix/linux systems, CronJob is a built-in service for performing regular
scheduled actions such as backups, report generation, and so on.

For this project, we are naively using kubernetes built-in cronJob service to
post a message to a third party API.

Since scheduled jobs differ from regular microservices hosted on a k8s cluster,
you may **not want** to use cronJobs to push data to your customers. You've been
warned ;)

## What does this proof-of-concept does?

Let's start by the small Go programm `email-service.go`

From the `main()` entrypoint, it will trigger the function `cronPublish(<my
message>)` and pass a message as input parameter.

The `cronPublish` function will format a short message containing a datetime
(now) and append the message supplied in the function call.

We make this `cronPublish`function POST an HTTP request to a third party API.
I had picked https://paste.c-net.org/ but any public pastebin site like
Sprunge will work the same way.

### how to (on a linux or mac machine)

1. have Docker running (follow instructions in README.md)
2. have minikube running (follow instructions from `minikube-setup.sh` file)
3. build a fresh Docker image **within your minikube environement** (follow
   same instructions as #2)
4. deploy the cronJob to your cluster using kubectl: `kubectl apply -f cronjob.yml`

### validate everything is working

```bash
kubectl get -n cron all

## expect
# > NAME                 SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
#   cron-email-service   */5 * * * *   False     0        <none>          67s
```

wait at least 5 minutes (job is scheduled to run every 5 minutes) and issue the
previous command a second time

```bash
## expect
# NAME                                    READY   STATUS    RESTARTS   AGE
# pod/cron-email-service-28027965-xh4wv   1/1     Running   0          2m18s

# NAME                               SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
# cronjob.batch/cron-email-service   */5 * * * *   False     1        2m18s           2m37s
```

use command `kubectl get pods -n cron` to list all cron service pods

pick the latest and checks its logs:

```bash
kubectl logs -n cron cron-email-service-28027965-xh4wv
```

## expect
# INFO: 2023/04/16 20:45:00 email-service.go:173: Starting up email-service
# INFO: 2023/04/16 20:45:00 email-service.go:175: trigger a cronmessage function
# INFO: 2023/04/16 20:45:00 email-service.go:57: now string is: 2023-04-16T20:45:00Z
# INFO: 2023/04/16 20:45:02 email-service.go:80: successfully posted cronmessage https://paste.c-net.org/SteeleSuppress
```

### before leaving for the day: delete your cronJob

`kubectl delete -f cronjob.yml`

and validate nothing is left running in your cluster: `kubectl get all -n cron`

expect: "No resources found in cron namespace."