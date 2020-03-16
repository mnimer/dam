

## Setup Instructions
1. Create PubSub Topic in Console
2. Create Subscription that calls the services/core/PubSubHandler
3. Register PubSub Notification listener
```Bash
    export topic = <Topic created in step #1>
```
```Bash
    export bucket = <Bucket Name to monitor> 
```
- register topic as a listener
```bash
gsutil notification create -t ${topic} -f json gs://${bucket}
```

4. Enable 'Google Cloud Storage JSON API'

## Deploy all Services
```bash
gcloud builds submit --config ./cloudbuild.yaml
```