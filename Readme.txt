Authentication
vscode ➜ /workspaces/Go Programming/techpuzzle7 $ gcloud auth activate-service-account go-dev@techpuzzle7.iam.gserviceaccount.com --key-file=sa-key.json
Activated service account credentials for: [go-dev@techpuzzle7.iam.gserviceaccount.com]

 41  gcloud init

   42  gcloud functions deploy techpuzzle7-function --entry-point HelloHTTP --runtime go116 --trigger-http --allow-unauthenticated

   43  curl -X POST https://us-central1-techpuzzle7.cloudfunctions.net/techpuzzle7-function -H "Content-Type:application/json" --data '{"name":"Keyboard Cat"}'


Storage Functions:
https://cloud.google.com/functions/docs/tutorials/storage

techpuzzle7_cloud_storage.go

gcloud functions deploy Techpuzzle7GCS \
--entry-point=Techpuzzle7GCS \
--runtime go116 \
--trigger-resource gs://techpuzzle7-invoices \
--trigger-event google.storage.object.finalize

Test:
touch gcf-test2.txt
gsutil cp gcf-test2.txt gs://techpuzzle7-invoices

gcloud functions logs read --limit 100

gcloud functions deploy process-invoice --runtime go116 --trigger-bucket techpuzzle7-invoices --entry-point ProcessInvoice 