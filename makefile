local-build:
	docker build -t mobydeck/ci-teams-notification .
	docker image prune -f


run: 
	export PLUGIN_WEBHOOKS='[{"webhook":"https://chat.googleapis.com/v1/spaces/AAAAgTnjGho/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=aXpTnIFH3w64wVPm8FOm5gfxoHPoBmevOeXpcIC1VJ0","provider":"google_chat","configs":{}}]' && \
	export PLUGIN_DEBUG=true && \
	export CI_COMMIT_SHA=1234567890 && \
	export CI_REPO=1234567890 && \
	export CI_COMMIT_AUTHOR=peter && \
	export CI_COMMIT_AUTHOR_AVATAR="https://cdn3.iconfinder.com/data/icons/round-default/64/add-1024.png" && \
	export CI_PREV_PIPELINE_STATUS=failure && \
	export CI_COMMIT_MESSAGE=1234567890 && \
	export CI_COMMIT_REF=1234567890 && \
	export CI_COMMIT_TAG=v1.0.1 && \
	export CI_PREV_COMMIT_URL="https://cdn3.iconfinder.com/data/icons/round-default/64/add-1024.png" && \
	export CI_PIPELINE_URL="https://cdn3.iconfinder.com/data/icons/round-default/64/add-1024.png" && \
	export CI_PIPELINE_FORGE_URL="https://cdn3.iconfinder.com/data/icons/round-default/64/add-1024.png" && \
	go run main.go
