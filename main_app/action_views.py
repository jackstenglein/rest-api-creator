from django.http import JsonResponse
from django.db.utils import IntegrityError
from django.views.decorators.csrf import csrf_exempt

from .models import Object, Action

import json



def create_action(request, project_id, object_id):
    """ Creates a new Action for the given user under the given Project and Object. """
    if not request.body:
        return JsonResponse({"error": "Missing request body"}, status=400)

    # Check parameters
    body_data = json.loads(request.body)
    if "name" not in body_data:
        return JsonResponse({"error": "Missing `name` parameter"}, status=400)
    if "route" not in body_data:
        return JsonResponse({"error": "Missing `route` parameter"}, status=400)
    if "method" not in body_data:
        return JsonResponse({"error": "Missing `method` parameter"}, status=400)

    try:
        object = Object.objects.get(pk=object_id, owner=request.user, project__id=project_id)
        new_action = Action(friendly_name=body_data["name"], kebab_case_name=body_data["name"], route=body_data["route"], method=body_data["method"], object=object)
        new_action.save()
        return JsonResponse({"message": "Action created"})
    except Object.DoesNotExist as e:
        return JsonResponse({"error": "Object id=" + str(object_id) + " does not exist"}, status=400)
    except KeyError as e:
        new_action.delete()
        return JsonResponse({"error": "Missing `" + e.args[0] + "` attribute parameter"}, status=400)
    except IntegrityError as e:
        if "action.name" in e.args[0]:
            error = "Action with this name already exists for this object"
        else:
            new_action.delete()
            error = "Unknown integrity error when creating action"
        return JsonResponse({"error": error}, status=400)

@csrf_exempt
def actions(request, project=None, object=None, action=None):
    if project == None:
        return JsonResponse({"error": "Invalid project"}, status=400)
    if object == None:
        return JsonResponse({"error": "Invalid object"}, status=400)
    if not request.user.is_authenticated:
        return JsonResponse({"error": "Not authenticated"}, status=401)
    # if request.method == "GET":
    #     return get_actions(request, project, object, action)
    if request.method == "POST":
        return create_action(request, project, object)

    return JsonResponse({"error": "Invalid HTTP method"}, status=405)
