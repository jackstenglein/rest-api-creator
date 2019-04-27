from django.shortcuts import render, redirect
from django.http import HttpResponse, JsonResponse
from django.forms import formset_factory
from django.contrib.auth import login, authenticate
from django.contrib.auth.forms import UserCreationForm

from .models import Object, Project, ProjectSerializer
from .forms import ObjectForm, AttributeForm
from django.views.decorators.csrf import csrf_exempt
import json


def create_project(request):
    """ Creates a new Project for the user in request.user. """
    body_data = json.loads(request.body)
    print(request.user)
    project = Project(owner=request.user, name=body_data["name"], language=1)
    project.save()
    return JsonResponse({"message": "Project created"}, status=201)


def get_projects(request):
    """ Gets all the Projects for the user in request.user. """
    projects = Project.objects.filter(owner=request.user)
    return JsonResponse({"projects": list(projects)}, JsonEncoder)


@csrf_exempt
def projects(request):
    """ Handles GET and POST requests on Projects. """
    if not request.user.is_authenticated:
        return JsonResponse({"error": "Not authenticated"}, status=401)
    if request.method == "GET":
        return get_projects(request)
    if request.method == "POST":
        return create_project(request)
    return JsonResponse({"error": "Invalid HTTP method"}, status=405)

