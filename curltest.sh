curl -X POST -H 'Content-Disposition: form-data; name="payload"' -H "Content-Type: application/json" -d @plex_hook_test.json "http://localhost:8080/plexhook"