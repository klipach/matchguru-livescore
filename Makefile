PROJECT_ID=match-guru-0iqc9r
FUNCTION_NAME=livescore
REGION=us-central1
PUBSUB_TOPIC=livescore-scheduler

deploy:
	gcloud functions deploy $(FUNCTION_NAME) \
	--region=$(REGION) \
	--project=$(PROJECT_ID) \
	--trigger-topic=$(PUBSUB_TOPIC) \
	--gen2 \
	--runtime=go123 \
	--timeout=15s \
	--max-instances=1 \
	--memory=128Mi \
	--entry-point=Sync \
	--source .

deploy_scheduler:
	gcloud scheduler jobs create pubsub $(PUBSUB_TOPIC) \
	--location=$(REGION) \
	--project=$(PROJECT_ID) \
	--schedule="every 1 minutes" \
	--time-zone="UTC" \
	--topic=$(PUBSUB_TOPIC) \
	--message-body="{}"

deploy_pubsub:
	gcloud pubsub topics create $(PUBSUB_TOPIC) \
	--project=$(PROJECT_ID) \
	--message-retention-duration=10m

get_function_url:
	gcloud functions describe $(FUNCTION_NAME) \
	--project=$(PROJECT_ID) \
	--format="value(url)"

firebase:
	firebase deploy --only hosting
