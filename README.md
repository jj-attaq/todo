curl http://localhost:8080/todo-list \--include \--header "Content-Type: application/json" \--request "POST" \--data '{"id": null, "item": "conquer the world", "finished": false, "uuid": ""}'

curl http://localhost:8080/todo-list/65467ef8-2f32-4805-b4fb-c55e74de1c68

curl http://localhost:8080/todo-list/update/4f9fd277-c17c-45a5-adb5-13f14f4f179a --header "Content-Type:application/json" \--request "POST"
