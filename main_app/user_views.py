from django.contrib.auth import login, authenticate
from django.contrib.auth.models import User
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
import json


@csrf_exempt
def signup(request):
    if request.method == "POST":
        print("Create new user")
        body_data = json.loads(request.body)
        print(body_data)
        user = User.objects.create_user(body_data["username"], body_data["email"], body_data["password"])
        user.save()
        login(request, user)
        return JsonResponse({"message": "User created"})


@csrf_exempt
def app_login(request):
    if request.method == "PUT":
        body_data = json.loads(request.body)

        if "username" not in body_data:
            return JsonResponse({"error": "Missing required `username` parameter"}, status=400)
        if "password" not in body_data:
            return JsonResponse({"error": "Missing required `password` parameter"}, status=400)

        user = authenticate(request, username=body_data["username"], password=body_data["password"])
        if user is not None:
            login(request, user)
            return JsonResponse({"message": "Logged in"})

        return JsonResponse({"error": "Incorrect username or password"}, status=400)

