from django.core.serializers.json import DjangoJSONEncoder
from .models import Attribute, AttributeSerializer, Object, ObjectSerializer, Project, ProjectSerializer


class JsonEncoder(DjangoJSONEncoder):
    def default(self, obj):
        if isinstance(obj, Attribute):
            return AttributeSerializer(obj).data
        if isinstance(obj, Project):
            return ProjectSerializer(obj).data
        if isinstance(obj, Object):
            return ObjectSerializer(obj).data
        return super().default(obj)
