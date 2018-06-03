from django.forms import ModelForm
from .models import Object, Attribute

class ObjectForm(ModelForm):
    class Meta:
        model = Object
        fields = ['name']

class AttributeForm(ModelForm):
    class Meta:
        model = Attribute
        fields = ['name', 'type', 'required']
