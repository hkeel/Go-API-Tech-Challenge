# ── Database ────────────────────────────────────────────────────────────────────

.PHONY: db_up
db_up:
	docker-compose up postgres

.PHONY: db_up_d
db_up_d:
	docker-compose up postgres -d

.PHONY: db_down
db_down:
	docker-compose down postgres

# ── API ─────────────────────────────────────────────────────────────────────────

.PHONY: run_app
run_app:
	docker-compose up

create-pr-description:
	git diff main...HEAD > differences.txt
	@if [ -z "${STORY_NUMBER}"]; then \
		echo "No JIRA story number found in branch name."; \
	else \
		echo "Please create a Pull Request description using the pull_request_template.md file attached with code changes from the differences.txt file attached. Generate two sentences max overview explaining the changes at a high level. Story number is ${STORY_NUMBER}. The output content must be in Markdown using source code format"