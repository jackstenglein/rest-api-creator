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
