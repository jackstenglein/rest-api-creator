from django.http import JsonResponse, FileResponse
from django.views.decorators.csrf import csrf_exempt
from django.core.files import File

from .models import Attribute, Object, Project

from .json_encoder import JsonEncoder
import json

import os
from zipfile import ZipFile
from rest_api_creator.settings import BASE_DIR



def create_project(request):
    """ Creates a new Project for the user in request.user. """
    body_data = json.loads(request.body)
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


def download_project(request, project_id):
    """ Returns the compiled Project. """
    if request.method != "GET":
        return JsonResponse({"error": "Invalid HTTP method"}, status=405)
    if not request.user.is_authenticated:
        return JsonResponse({"error": "Not authenticated"}, status=401)

    try:
        project = Project.objects.get(pk=project_id, owner=request.user)
        objects = Object.objects.filter(owner=request.user, project=project)

        dirpath = get_project_dir(request.user.username, project_id)
        filepaths = []
        try:
            os.makedirs(dirpath)
        except FileExistsError:
            pass

        for object in objects:
            # Create a Python file object using open() and the with statement
            filepath = os.path.join(dirpath, object.name + ".js")
            filepaths.append(filepath)
            with open(filepath, 'w') as f:
                myfile = File(f)
                myfile.write("// api/models/" + object.name + ".js\n\n")
                myfile.write("module.exports = {\n    attributes: {")

                attributes = Attribute.objects.filter(object=object)
                for attribute in attributes:
                    myfile.write("\n        " + attribute.name + ": { type: '")
                    myfile.write(Attribute.DATA_TYPE_TO_SAILS_TYPE[attribute.type])
                    myfile.write("', required: " + str(attribute.required).lower() + " },")

                myfile.write("\n    }\n}\n")

        zip_path = dirpath + "/" + project.name + ".zip"
        with ZipFile(zip_path,'w') as zip:
            # writing each file one by one
            for file in filepaths:
                file_name = file[file.rindex("/")+1 : ]
                zip.write(file, arcname=project.name + "/api/models/" + file_name)

        return FileResponse(open(zip_path, 'rb'), as_attachment=True)
    except Project.DoesNotExist:
        return JsonResponse({"error": "Project id=" + str(project_id) + " does not exist"}, status=400)


def get_project_dir(username, project_id):
    return os.path.join(BASE_DIR, "compiled_projects", username, "projects", str(project_id), "objects")



