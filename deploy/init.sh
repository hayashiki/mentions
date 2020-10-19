#!/usr/bin/env bash

#This scripts is to be run once per environments
# * enables the required APIs
# * create a service account and set the appropriate right to run the cloud functions
# * create the cloud tasks queues
#

echo "initializing GCP Project..."

#Cloud Functions
echo "gcloud services enable cloudfunctions.googleapis.com"
gcloud services enable cloudfunctions.googleapis.com

#Cloud Tasks
echo "gcloud services enable cloudtasks.googleapis.com"
gcloud services enable cloudtasks.googleapis.com

#GAE
echo "gcloud services enable appengine.googleapis.com"
gcloud services enable cloudbuild.googleapis.com

#GAE
echo "gcloud services enable appengine.googleapis.com"
gcloud services enable appengine.googleapis.com

echo "Creating deploy service account"
gcloud iam service-accounts create github-deploy-action \
    --description="Service Account for the github deploy action" \
    --display-name="Deploy Action SA"

# cloudbuild
# appengine
# cloudstorage
echo "Add role to service account: roles/appengine.appViewer"
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member "serviceAccount:github-deploy-action@${PROJECT_ID}.iam.gserviceaccount.com" \
  --roles/appengine.appAdmin \
  --roles/cloudbuild.builds.builder \
  --roles/storage.admin


