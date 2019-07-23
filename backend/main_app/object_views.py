from django.http import JsonResponse
from django.db.utils import IntegrityError
from django.views.decorators.csrf import csrf_exempt
from django.forms.models import model_to_dict

from .models import Object, Project, ProjectSerializer
from .json_encoder import JsonEncoder

import json


def get_objects(request, project_id, object_id=None):
    """ Gets the specified objects for the authenticated user under the given project. """
    if object_id is not None:
        # Get single
        try:
            object = Object.objects.get(pk=object_id, owner=request.user, project__id=project_id)
            attributes = object.attribute_set.all()
            actions = object.action_set.all()
            object = model_to_dict(object)
            object["attributes"] = list(attributes)
            object["actions"] = list(actions)
            return JsonResponse({"object": object}, JsonEncoder)
        except Object.DoesNotExist as e:
            return JsonResponse({"error": "Object id=" + str(object_id) + " does not exist"}, status=400)
    else:
        # Get many
        objects = Object.objects.filter(owner=request.user, project__id=project_id)
        return JsonResponse({"objects": list(objects)}, JsonEncoder)


def create_object(request, project_id):
    """ Creates a new Object for the given user under the given Project. """
    if not request.body:
        return JsonResponse({"error": "Missing request body"}, status=400)

    body_data = json.loads(request.body)
    if "name" not in body_data:
        return JsonResponse({"error": "Missing `name` parameter"}, status=400)

    try:
        project = Project.objects.get(pk=project_id, owner=request.user)
        new_object = Object(name=body_data["name"], project=project, owner=request.user)
        new_object.save()
        if "attributes" in body_data:
            for attribute in body_data["attributes"]:
                new_object.attribute_set.create(
                    name=attribute["name"],
                    required=attribute.get("required", False),
                    type=attribute["type"]
                )
        print(new_object)
        return JsonResponse({"message": "Object created"})
    except Project.DoesNotExist as e:
        return JsonResponse({"error": "Project id=" + str(project_id) + " does not exist"}, status=400)
    except KeyError as e:
        new_object.delete()
        return JsonResponse({"error": "Missing `" + e.args[0] + "` attribute parameter"}, status=400)
    except IntegrityError as e:
        if "object.name" in e.args[0]:
            error = "Object with this name already exists for this project"
        else:
            new_object.delete()
            error = "Duplicate attribute names for this object"
        return JsonResponse({"error": error}, status=400)

@csrf_exempt
def objects(request, project=None, object=None):    
    if request.method == "OPTIONS":
        return HttpResponse()
    if project == None:
        return JsonResponse({"error": "Invalid project"}, status=400)
    if not request.user.is_authenticated:
        return JsonResponse({"error": "Not authenticated"}, status=401)
    if request.method == "GET":
        return get_objects(request, project, object)
    if request.method == "POST":
        return create_object(request, project)

    return JsonResponse({"error": "Invalid HTTP method"}, status=405)
