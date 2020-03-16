## Deploy
Deploy your workflow (in the first deployment you will be prompted to enable API):
``` bash
gcloud alpha workflows deploy testworkflow --source=test.yaml
```

## Execute
Execute your workflow (in the first execution you will be prompted to enable API):
```bash
gcloud alpha workflows execute testworkflow
```
This command will also echo a command you can copy/paste & run to see its execution status and results. 
When you run it, Result field will hold a workflow output with a Wikipedia extract for a given day. 